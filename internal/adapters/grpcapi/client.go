package grpcapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/AcroManiac/micropic/internal/domain/interfaces"

	"google.golang.org/grpc"
)

type CacheClientImpl struct {
	conn   *grpc.ClientConn
	client CacheClient
}

func NewCacheClientImpl(host string, port int) interfaces.CacheClient {
	// Start gRPC client
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure())
	if err != nil {
		logger.Fatal("could not connect gRPC server", "error", err)
	}

	// Create gRPC cache client object
	client := NewCacheClient(conn)
	return &CacheClientImpl{
		conn:   conn,
		client: client,
	}
}

// Save preview to cache via gRPC
func (c *CacheClientImpl) Save(preview *entities.Preview) error {
	// Create timed context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send gRPC request
	_, err := c.client.SavePreview(ctx, &SavePreviewRequest{
		Preview: convertPreviewToProtobuf(preview),
	})
	if err != nil {
		logger.Error("error sending preview to cache", "error", err)
		return err
	}
	return nil
}

// Get preview from cache via gRPC
func (c *CacheClientImpl) Get(params *entities.PreviewParams) (*entities.Preview, error) {
	// Create timed context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send gRPC request
	response, err := c.client.GetPreview(ctx, &GetPreviewRequest{
		Url:    params.URL,
		Width:  uint32(params.Width),
		Height: uint32(params.Height),
	})
	if err != nil {
		logger.Error("error getting preview from cache", "error", err)
		return nil, err
	}

	// Check response for errors
	errText := response.GetError()
	if errText != "" {
		logger.Error("cache returned error", "error", errText)
		return nil, errors.New(errText)
	}

	// Get preview from cache response
	preview := convertProtobufToPreview(response.GetPreview())
	return preview, nil
}
