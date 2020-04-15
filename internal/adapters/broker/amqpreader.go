package broker

import (
	"context"
	"io"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type AmqpReader struct {
	ctx  context.Context
	cwq  *ChannelWithQueue
	msgs <-chan amqp.Delivery
}

func NewAmqpReader(ctx context.Context, conn *amqp.Connection) io.ReadCloser {
	// Create amqp channel and queue
	queueName := queueName
	ch, err := NewChannelWithQueue(conn, &queueName)
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
		true,        // exclusive
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

// Read one message from RabbitMQ queue. Returns message length in bytes
func (r *AmqpReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		logger.Debug("Context cancelled", "caller", "AmqpReader")
	case message, ok := <-r.msgs:
		if ok {
			n = copy(p, message.Body)
		}
	}
	return
}

func (r *AmqpReader) Close() error {
	if err := r.cwq.Close(); err != nil {
		return errors.Wrap(err, "failed closing gateway output channel")
	}
	return nil
}
