package broker

import (
	"encoding/json"

	"github.com/ahamtat/micropic/internal/adapters/logger"
	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpWriter struct {
	cwq        *ChannelWithQueue
	routingKey string
}

func NewAmqpWriter(conn *amqp.Connection, queueName, routingKey string) *AmqpWriter {
	// Create amqp channel and queue
	ch, err := NewChannelWithQueue(conn, queueName, routingKey)
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
			Type:          env.Metadata.Type,
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
