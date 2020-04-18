package main

import (
	"github.com/AcroManiac/micropic/internal/adapters/broker"
	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/AcroManiac/micropic/internal/domain/entities"
)

type RMQSender struct {
	w *broker.AmqpWriter
}

func NewRMQSender(w *broker.AmqpWriter) *RMQSender {
	return &RMQSender{w: w}
}

func (s *RMQSender) Send(response *entities.Response, obj interface{}) {
	// Get correlation
	if obj == nil {
		logger.Error("no correlation id in params")
		return
	}
	correlationID, ok := obj.(string)
	if !ok {
		logger.Error("error type asserting correlation id to string")
	}

	// Send preview response to RabbitMQ
	env := broker.CreateEnvelope(response, correlationID,
		entities.MessageTypeToString(entities.PreviewerResponse))
	if err := s.w.WriteEnvelope(env); err != nil {
		logger.Error("error writing envelope to RabbitMQ", "error", err, "caller", "Send")
	}
}
