package broker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/cenkalti/backoff/v4"
)

const (
	exchangeName = "micropic"
	queueName    = "previewerQueue"
	routingKey   = "previewerQueue.images"
)

type Manager struct {
	connURL string
	Conn    *amqp.Connection
	Done    chan error
	wr      io.WriteCloser
	rd      io.ReadCloser
}

func NewManager(protocol, user, password, host string, port int) *Manager {
	// Create manager object
	m := &Manager{
		connURL: fmt.Sprintf("%s://%s:%s@%s:%d/", protocol, user, password, host, port),
		Conn:    nil,
		Done:    make(chan error),
	}

	// Open connection to broker
	if err := m.connect(); err != nil {
		logger.Error("RabbitMQ connection failed", "error", err)
		return nil
	}

	return m
}

func (m *Manager) connect() (err error) {
	// Open RabbitMQ connection
	m.Conn, err = amqp.Dial(m.connURL)
	if err != nil {
		return
	}

	return
}

func (m *Manager) ConnectionListener(ctx context.Context) {
	select {
	case <-ctx.Done():
		break
	case connErr := <-m.Conn.NotifyClose(make(chan *amqp.Error)):
		logger.Error("RabbitMQ connection is closed", "error", connErr.Error())
		// Notify clients
		m.Done <- errors.New("connection closed")
	}
}

func (m *Manager) Reconnect(ctx context.Context) error {
	// Close i/o channels
	m.closeIOChannels()

	// Create reconnect backoff
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	// Do reconnect loop
	boCtx := backoff.WithContext(be, context.Background())
	for {
		boTime := boCtx.NextBackOff()
		if boTime == backoff.Stop {
			return errors.New("backoff timer elapsed")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(boTime):
			if err := m.connect(); err != nil {
				logger.Error("couldn't reconnect", "error", err)
				continue
			}
			logger.Info("Reconnect to RabbitMQ succeeded")
			return nil
		}
	}
}

func (m *Manager) GetWriter() io.Writer {
	if m.wr == nil {
		// Create broker writer
		m.wr = NewAmqpWriter(m.Conn)
	}
	return m.wr
}

func (m *Manager) GetReader(ctx context.Context) io.Reader {
	if m.rd == nil {
		// Create broker reader
		m.rd = NewAmqpReader(ctx, m.Conn)
	}
	return m.rd
}

func (m *Manager) closeIOChannels() {
	if m.rd != nil {
		m.rd.Close()
		m.rd = nil
	}
	if m.wr != nil {
		m.wr.Close()
		m.wr = nil
	}
}

func (m *Manager) Close() error {
	// Close i/o channels
	m.closeIOChannels()

	// Close connection notify channel
	if m.Done != nil {
		close(m.Done)
	}

	// Close connection
	if m.Conn != nil {
		if err := m.Conn.Close(); err != nil {
			logger.Error("failed closing connection", "error", err)
		}
	}

	return nil
}
