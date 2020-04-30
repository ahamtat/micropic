package broker

import (
	"context"
	"encoding/json"

	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/AcroManiac/micropic/internal/domain/interfaces"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpReader struct {
	ctx  context.Context
	cwq  *ChannelWithQueue
	msgs <-chan amqp.Delivery
}

func NewAmqpReader(ctx context.Context, conn *amqp.Connection, queueName, routingKey string) *AmqpReader {
	// Create amqp channel and queue
	ch, err := NewChannelWithQueue(conn, queueName, routingKey)
	if err != nil {
		logger.Error("failed creating amqp channel and queue",
			"error", err, "queue", queueName,
			"caller", "NewAmqpReader")
		return nil
	}

	// Create consuming channel
	msgs, err := ch.Ch.Consume(
		ch.Que.Name, // queue
		"",          // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	if err != nil {
		logger.Error("failed to register a consumer",
			"error", err, "queue", queueName,
			"caller", "NewAmqpReader")
		return nil
	}

	// Return reader object
	return &AmqpReader{
		ctx:  ctx,
		cwq:  ch,
		msgs: msgs,
	}
}

// ReadEnvelope reads and unmarshals message from RabbitMQ queue. Returns message envelope or error
func (r *AmqpReader) ReadEnvelope() (env *AmqpEnvelope, close bool, err error) {
	select {
	case <-r.ctx.Done():
		logger.Debug("Context cancelled", "caller", "ReadEnvelope")
		close = true
	case message, ok := <-r.msgs:
		if ok {
			// Create message buffer
			bodyLength := len(message.Body)
			buffer := make([]byte, bodyLength)
			n := copy(buffer, message.Body)
			if n != bodyLength {
				err = errors.Wrap(err, "error copying message body to buffer")
				return
			}

			// Create empty message object
			var inputMessage interfaces.Message
			messageType, ok := entities.StringToMessageType(message.Type)
			if !ok {
				err = errors.New("error converting message type")
				return
			}
			switch messageType {
			case entities.ProxyingRequest:
				inputMessage = &entities.Request{}
			case entities.PreviewerResponse:
				inputMessage = &entities.Response{}
			}

			// Unmarshal input message from JSON to structure
			err = json.Unmarshal(buffer, &inputMessage)
			if err != nil {
				err = errors.Wrap(err, "can not unmarshal incoming gateway message")
				return
			}

			// Print copy of incoming message to log
			r.PrintMessage(inputMessage)

			// Create envelope
			env = &AmqpEnvelope{
				Message: inputMessage,
				Metadata: &AmqpMetadata{
					CorrelationID: message.CorrelationId,
					Type:          message.Type,
				},
			}
		}
	}
	return
}

// PrintMessage prints incoming message to log
func (r AmqpReader) PrintMessage(message interfaces.Message) {
	// Slim long preview
	response, ok := message.(*entities.Response)
	if ok {
		respCopy := &entities.Response{
			Preview: &entities.Preview{
				Params: response.Preview.Params,
				Image: func(p []byte) []byte {
					if p != nil {
						return []byte("Some Base64 code ;)")
					}
					return nil
				}(response.Preview.Image),
			},
			Status: response.Status,
		}
		logger.Debug("Response from previewer", "response", respCopy)
		return
	}
	logger.Debug("Request from HTTP proxy", "request", message)
}

// Close function releases RabbitMQ channel and corresponding queue
func (r *AmqpReader) Close() error {
	if err := r.cwq.Close(); err != nil {
		return errors.Wrap(err, "failed closing gateway output channel")
	}
	return nil
}
