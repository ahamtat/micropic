package main

import (
	"context"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"
	"github.com/AcroManiac/micropic/internal/domain/usecases"

	"github.com/AcroManiac/micropic/internal/domain/entities"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/AcroManiac/micropic/internal/adapters/broker"
	"github.com/streadway/amqp"
)

type RMQListener struct {
	ctx       context.Context
	cancel    context.CancelFunc
	in        *broker.AmqpReader
	out       *broker.AmqpWriter
	previewer *usecases.Previewer
}

func NewRMQListener(conn *amqp.Connection, quality int) *RMQListener {
	// Create cancel context
	ctx, cancel := context.WithCancel(context.Background())

	in := broker.NewAmqpReader(ctx, conn, broker.RequestQueueName, broker.RequestRoutingKey)
	out := broker.NewAmqpWriter(conn, broker.ResponseQueueName, broker.ResponseRoutingKey)

	previewer := usecases.NewPreviewer(
		NewHTTPImageLoader(),
		NewImageProcessor(quality),
		[]interfaces.Sender{
			NewRMQSender(out),
			NewCacheSender(),
		})

	return &RMQListener{
		ctx:       ctx,
		cancel:    cancel,
		in:        in,
		out:       out,
		previewer: previewer,
	}
}

// Close reading and writing channels
func (l *RMQListener) Close() {
	l.Stop()

	// Close i/o channels
	if err := l.in.Close(); err != nil {
		logger.Error("error closing RabbitMQ reader", "error", err)
	}
	if err := l.out.Close(); err != nil {
		logger.Error("error closing RabbitMQ writer", "error", err)
	}
}

// Start message receiving and processing
func (l *RMQListener) Start() {
	// Read and process messages from previewer
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			// Read input message
			envelope, toBeClosed, err := l.in.ReadEnvelope()
			if err != nil {
				logger.Error("error reading input channel", "error", err)
				break
			}
			if toBeClosed {
				// Reading channel is to be closed possibly
				break
			}

			// Check for proxy requests
			if len(envelope.Metadata.CorrelationID) > 0 {
				// Get request data
				request, ok := envelope.Message.(*entities.Request)
				if !ok || request == nil {
					logger.Error("wrong proxy request data")
					break
				}

				// WARNING!!! CRIPPLED CODE!!!
				// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				l.previewer.Process(request, envelope.Metadata.CorrelationID)
				// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				// For tutorial purposes this code was made non-concurrent
				// In real application one should use:
				// go l.previewer.Process(request, envelope.Metadata.CorrelationID)
			}
		}
	}
}

// Stop message listening
func (l *RMQListener) Stop() {
	// Fire context cancelling
	l.cancel()
}
