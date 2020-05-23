package grpcapi

import "github.com/ahamtat/micropic/internal/domain/entities"

func convertParamsToProtobuf(params *entities.PreviewParams) *PreviewParams {
	return &PreviewParams{
		Url:    params.URL,
		Width:  uint32(params.Width),
		Height: uint32(params.Height),
	}
}

func convertPreviewToProtobuf(preview *entities.Preview) *Preview {
	return &Preview{
		Params: convertParamsToProtobuf(preview.Params),
		Image:  preview.Image,
	}
}

func convertProtobufToPreview(protobuf *Preview) *entities.Preview {
	return &entities.Preview{
		Params: &entities.PreviewParams{
			Width:  int(protobuf.Params.Width),
			Height: int(protobuf.Params.Height),
			URL:    protobuf.Params.Url,
		},
		Image: protobuf.Image,
	}
}
