package grpcapi

import (
	context "context"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/pkg/errors"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"
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
	logger.Debug("Incoming SavePreview request", "request", request.String())

	// Check input params
	if request == nil {
		err := errors.New("no input params")
		logger.Error(err.Error())
		return &SavePreviewResponse{Error: err.Error()}, err
	}

	// Save preview in cache
	err := s.cache.Save(&entities.Preview{
		Params: &entities.PreviewParams{
			Width:  int(request.Preview.Params.Width),
			Height: int(request.Preview.Params.Height),
			URL:    request.Preview.Params.Url,
		},
		Image: request.Preview.Image,
	})
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
			Preview: &Preview{
				Params: &PreviewParams{
					Url:    preview.Params.URL,
					Width:  uint32(preview.Params.Width),
					Height: uint32(preview.Params.Height),
				},
				Image: preview.Image,
			},
		},
	}, nil
}
