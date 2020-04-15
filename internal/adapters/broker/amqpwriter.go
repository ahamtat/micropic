package broker

import (
	"encoding/json"
	"io"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpWriter struct {
	cwq        *ChannelWithQueue
	routingKey string
}

func NewAmqpWriter(conn *amqp.Connection) io.WriteCloser {
	// Create amqp channel and queue
	ch, err := NewChannelWithQueue(conn, nil)
	if err != nil {
		logger.Error("failed creating amqp channel and queue",
			"error", err,
			"caller", "NewAmqpWriter")
		return nil
	}

	return &AmqpWriter{
		cwq:        ch,
		routingKey: routingKey,
	}
}

// Write message to RabbitMQ broker.
// Returns message length on success or error if any
func (w *AmqpWriter) Write(p []byte) (n int, err error) {
	if w.cwq.Ch == nil {
		return 0, errors.New("no output channel defined")
	}

	// Send message to gateway
	err = w.cwq.Ch.Publish(
		exchangeName, // exchange
		w.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        p,
		})
	if err != nil {
		return 0, errors.Wrap(err, "failed to publish a message")
	}

	return len(p), nil
}

// WriteEnvelope sends AMQP envelope to RabbitMQ broker, returns error object or nil
func (w *AmqpWriter) WriteEnvelope(env *AmqpEnvelope) error {
	if w.cwq.Ch == nil {
		return errors.New("no output channel defined")
	}

	// Marshall message to JSON
	buffer, err := json.Marshal(env.Message)
	if err != nil {
		return errors.Wrap(err, "error marshalling RPC request to JSON")
	}

	// Send message with metadata to gateway queue
	err = w.cwq.Ch.Publish(
		exchangeName, // exchange
		w.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: env.Metadata.CorrelationID,
			ReplyTo:       env.Metadata.ReplyTo,
			Body:          buffer,
		})
	if err != nil {
		return errors.Wrap(err, "failed to publish a message")
	}

	return nil
}

func (w *AmqpWriter) Close() error {
	if err := w.cwq.Close(); err != nil {
		return errors.Wrap(err, "failed closing writer channel")
	}
	return nil
}
