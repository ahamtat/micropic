package grpcapi

import (
	context "context"

	"github.com/ahamtat/micropic/internal/adapters/logger"
	"github.com/ahamtat/micropic/internal/domain/entities"
	"github.com/pkg/errors"

	"github.com/ahamtat/micropic/internal/domain/interfaces"
)

// CacheServerImpl implementation for gRPC API
type CacheServerImpl struct {
	cache interfaces.Cache
}

// NewCacheServerImpl constructor
func NewCacheServerImpl(cache interfaces.Cache) CacheServer {
	return &CacheServerImpl{cache: cache}
}

// SavePreview gRPC call handler saves preview in the internal cache
func (s *CacheServerImpl) SavePreview(ctx context.Context, request *SavePreviewRequest) (*SavePreviewResponse, error) {
	logger.Debug("Incoming SavePreview request", "request", request.Preview.Params)

	// Check input params
	if request == nil {
		err := errors.New("no input params")
		logger.Error(err.Error())
		return &SavePreviewResponse{Error: err.Error()}, err
	}

	// Save preview in cache
	err := s.cache.Save(convertProtobufToPreview(request.Preview))
	if err != nil {
		err = errors.Wrap(err, "error saving preview in cache")
		logger.Error(err.Error())
		return &SavePreviewResponse{Error: err.Error()}, err
	}

	logger.Debug("SavePreview processed successfully")
	return &SavePreviewResponse{Error: ""}, nil
}

// GetPreview gRPC call handler returns preview from internal cache
func (s *CacheServerImpl) GetPreview(ctx context.Context, request *GetPreviewRequest) (*GetPreviewResponse, error) {
	logger.Debug("Incoming GetPreview request", "request", request.String())

	// Check input params
	if request == nil {
		err := errors.New("no input params")
		logger.Error(err.Error())
		return &GetPreviewResponse{
			Result: &GetPreviewResponse_Error{Error: err.Error()},
		}, err
	}

	// Get preview from cache
	preview, err := s.cache.Get(&entities.PreviewParams{
		Width:  int(request.Width),
		Height: int(request.Height),
		URL:    request.Url,
	})
	if err != nil {
		err = errors.Wrap(err, "failed getting preview from cache")
		logger.Error(err.Error())
		return &GetPreviewResponse{
			Result: &GetPreviewResponse_Error{Error: err.Error()},
		}, err
	}

	logger.Debug("GetPreview returned preview")
	return &GetPreviewResponse{
		Result: &GetPreviewResponse_Preview{
			Preview: convertPreviewToProtobuf(preview),
		},
	}, nil
}
