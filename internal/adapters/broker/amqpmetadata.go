package broker

import (
	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/google/uuid"
)

// AmqpMetadata holds extra data for AMQP message
type AmqpMetadata struct {
	CorrelationID string
	ReplyTo       string
}

// CreateCorrelationID returns correlation UUID
func CreateCorrelationID() string {
	return uuid.New().String()
}

// AmqpEnvelope structure to hold IoT message with AMQP metadata
type AmqpEnvelope struct {
	Message  *entities.Image
	Metadata *AmqpMetadata
}
