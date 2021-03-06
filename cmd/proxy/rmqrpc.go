package main

import (
	"context"
	"sync"

	"github.com/ahamtat/micropic/internal/adapters/broker"

	"github.com/ahamtat/micropic/internal/domain/entities"

	"github.com/pkg/errors"

	"github.com/ahamtat/micropic/internal/domain/interfaces"

	"github.com/ahamtat/micropic/internal/adapters/logger"
	"github.com/streadway/amqp"
)

// RMQRPC holds objects for making Remote Procedure Calls via RabbitMQ broker
// See https://medium.com/@eugenfedchenko/rpc-over-rabbitmq-golang-ff3d2b312a69
// && https://www.rabbitmq.com/tutorials/tutorial-six-go.html
type RMQRPC struct {
	ctx      context.Context
	cancel   context.CancelFunc
	in       *broker.AmqpReader
	out      *broker.AmqpWriter
	rpcMx    sync.Mutex
	rpcCalls rpcPendingCallMap
}

type rpcPendingCall struct {
	done chan bool
	data interfaces.Message
}

type rpcPendingCallMap map[string]*rpcPendingCall

// NewRPC constructor
func NewRPC(conn *amqp.Connection) *RMQRPC {
	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	in := broker.NewAmqpReader(ctx, conn, broker.ResponseQueueName, broker.ResponseRoutingKey)
	out := broker.NewAmqpWriter(conn, broker.RequestQueueName, broker.RequestRoutingKey)

	return &RMQRPC{
		ctx:      ctx,
		cancel:   cancel,
		in:       in,
		out:      out,
		rpcMx:    sync.Mutex{},
		rpcCalls: make(rpcPendingCallMap),
	}
}

// Close reading and writing channels
func (rpc *RMQRPC) Close() {
	rpc.Stop()

	// Close pending calls to quit blocked goroutines
	rpc.rpcMx.Lock()
	for _, call := range rpc.rpcCalls {
		close(call.done)
	}
	rpc.rpcMx.Unlock()

	// Close i/o channels
	if err := rpc.in.Close(); err != nil {
		logger.Error("error closing RabbitMQ reader", "error", err)
	}
	if err := rpc.out.Close(); err != nil {
		logger.Error("error closing RabbitMQ writer", "error", err)
	}
}

// Start functions make separate goroutine for message receiving and processing
func (rpc *RMQRPC) Start() {
	// Read and process messages from previewer
	for {
		select {
		case <-rpc.ctx.Done():
			return
		default:
			// Read input message
			inputEnvelope, toBeClosed, err := rpc.in.ReadEnvelope()
			if err != nil {
				logger.Error("error reading channel", "error", err)
				break
			}
			if toBeClosed {
				// Reading channel possibly is to be closed
				break
			}

			// Check for RMQRPC responses
			if len(inputEnvelope.Metadata.CorrelationID) > 0 {
				// Make pending call
				rpc.rpcMx.Lock()
				rpcCall, ok := rpc.rpcCalls[inputEnvelope.Metadata.CorrelationID]
				rpc.rpcMx.Unlock()
				if ok {
					rpcCall.data = inputEnvelope.Message
					rpcCall.done <- true
				}
				break
			}
		}
	}
}

// Stop message processing and writing off status to database
func (rpc *RMQRPC) Stop() {
	// Stop goroutines - fire context cancelling
	rpc.cancel()
}

// SendRequest sends tasks for previewer via RabbitMQ broker and
// blocks execution until response or timeout
func (rpc *RMQRPC) SendRequest(ctx context.Context, request *entities.Request) (response *entities.Response, err error) {
	// Create message envelope and write it to broker
	correlationID := broker.CreateCorrelationID()
	env := broker.CreateEnvelope(request, correlationID,
		entities.MessageTypeToString(entities.ProxyingRequest))
	err = rpc.out.WriteEnvelope(env)
	if err != nil {
		return nil, errors.Wrap(err, "error writing envelope to RabbitMQ")
	}

	// Create and keep pending object
	rpcCall := &rpcPendingCall{done: make(chan bool)}
	rpc.rpcMx.Lock()
	rpc.rpcCalls[correlationID] = rpcCall
	rpc.rpcMx.Unlock()

	// Wait until response comes or timeout
	select {
	case <-rpcCall.done:
		response, _ = rpcCall.data.(*entities.Response)
	case <-ctx.Done():
		err = errors.New("timeout elapsed on RMQRPC request sending")
	}

	// Release pending object
	rpc.rpcMx.Lock()
	delete(rpc.rpcCalls, correlationID)
	rpc.rpcMx.Unlock()

	// Return response to caller
	return
}
