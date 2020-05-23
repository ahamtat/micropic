package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/muesli/smartcrop/nfnt"

	"github.com/muesli/smartcrop"

	"github.com/ahamtat/micropic/internal/adapters/logger"

	"github.com/ahamtat/micropic/internal/domain/entities"
)

// ImageProcessor
type ImageProcessor struct {
	quality int
}

// NewImageProcessor constructor
func NewImageProcessor(quality int) *ImageProcessor {
	return &ImageProcessor{quality: quality}
}

func errorResponse(request *entities.Request, errText string) *entities.Response {
	return &entities.Response{
		Preview: &entities.Preview{
			Params: request.Params,
			Image:  []byte{},
		},
		Status: entities.Status{
			Code: http.StatusInternalServerError,
			Text: errText,
		},
	}
}

type subImager interface {
	SubImage(image.Rectangle) image.Image
}

// Process source image with params in request and
// return preview in Base64 format or error status
func (p *ImageProcessor) Process(srcImage []byte, request *entities.Request) *entities.Response {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(srcImage))
	if err != nil {
		logger.Error("error decoding source image", "error", err)
		return errorResponse(request, err.Error())
	}
	logger.Debug("Source image decoded successfully", "format", format)

	// Make preview cropping from decoded image
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	cropArea, err := analyzer.FindBestCrop(img, request.Params.Width, request.Params.Height)
	if err != nil {
		logger.Error("failed searching best crop area", "error", err)
		return errorResponse(request, err.Error())
	}
	logger.Debug("Best crop", "area", cropArea)

	// Crop image with requested aspect ratio
	si, ok := img.(subImager)
	if !ok {
		errText := "failed cropping preview subimage"
		logger.Error(errText)
		return errorResponse(request, errText)
	}
	croppedImage := si.SubImage(cropArea)
	logger.Debug("Cropped image dimensions",
		"width", croppedImage.Bounds().Dx(), "height", croppedImage.Bounds().Dy())

	// Resize image to fit requested params
	resizedImage := resizer.Resize(croppedImage, uint(request.Params.Width), uint(request.Params.Height))

	// In-memory buffer to store JPEG image before we Base64 encode it
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	if err := jpeg.Encode(&buff, resizedImage, &jpeg.Options{Quality: p.quality}); err != nil {
		logger.Error("failed encoding preview to JPEG", "error", err)
		return errorResponse(request, err.Error())
	}

	// Encode the bytes in the buffer to a base64 string
	preview := make([]byte, base64.StdEncoding.EncodedLen(len(buff.Bytes())))
	base64.StdEncoding.Encode(preview, buff.Bytes())
	logger.Debug("Preview made successfully")

	// Preview made successfully
	return &entities.Response{
		Preview: &entities.Preview{
			Params: request.Params,
			Image:  preview,
		},
		Status: entities.Status{
			Code: http.StatusOK,
			Text: http.StatusText(http.StatusOK),
		},
	}
}
