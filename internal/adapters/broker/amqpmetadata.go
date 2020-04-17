package broker

import (
	"github.com/AcroManiac/micropic/internal/domain/interfaces"
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
