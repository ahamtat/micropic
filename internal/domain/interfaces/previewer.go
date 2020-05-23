package interfaces

import "github.com/ahamtat/micropic/internal/domain/entities"

// ImageLoader interface
type ImageLoader interface {
	// Load image from external source containing in request
	Load(request *entities.Request) ([]byte, *entities.Status)
}

// ImageProcessor interface
type ImageProcessor interface {
	// Process source image with params from request and return preview response
	Process(srcImage []byte, request *entities.Request) *entities.Response
}

// Sender interface
type Sender interface {
	// Send response to consumer
	Send(response *entities.Response, obj interface{})
}
