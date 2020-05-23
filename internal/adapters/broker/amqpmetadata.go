package broker

import (
	"github.com/ahamtat/micropic/internal/domain/interfaces"
	"github.com/google/uuid"
)

// AmqpMetadata holds extra data for AMQP message
type AmqpMetadata struct {
	CorrelationID string
	Type          string
}

// CreateCorrelationID returns correlation UUID
func CreateCorrelationID() string {
	return uuid.New().String()
}

// AmqpEnvelope holds message with AMQP metadata
type AmqpEnvelope struct {
	Message  interfaces.Message
	Metadata *AmqpMetadata
}

// CreateEnvelope creates and fills message envelope
func CreateEnvelope(message interfaces.Message, correlationID, dataType string) *AmqpEnvelope {
	return &AmqpEnvelope{
		Message: message,
		Metadata: &AmqpMetadata{
			CorrelationID: correlationID,
			Type:          dataType,
		},
	}
}
