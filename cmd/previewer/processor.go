package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/muesli/smartcrop/nfnt"

	"github.com/muesli/smartcrop"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/AcroManiac/micropic/internal/domain/entities"
)

// ImageProcessor
type ImageProcessor struct {
	quality int
}

// NewImageProcessor constructor
func NewImageProcessor(quality int) *ImageProcessor {
	return &ImageProcessor{quality: quality}
}

func errorResponse(filename, errText string) *entities.Response {
	return &entities.Response{
		Preview:  nil,
		Filename: filename,
		Status: entities.Status{
			Code: http.StatusInternalServerError,
			Text: errText,
		},
	}
}

// Process source image with params in request and
// return preview in Base64 format or error status
func (p *ImageProcessor) Process(srcImage []byte, request *entities.Request) *entities.Response {
	// Extract filename
	filename := request.URL[strings.LastIndex(request.URL, "/")+1:]

	// Decode image
	img, format, err := image.Decode(bytes.NewReader(srcImage))
	if err != nil {
		logger.Error("error decoding source image", "error", err)
		return errorResponse(filename, err.Error())
	}
	logger.Debug("Source image decoded successfully", "filename", filename, "format", format)

	// Make preview cropping from decoded image
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	cropArea, err := analyzer.FindBestCrop(img, request.Width, request.Height)
	if err != nil {
		logger.Error("failed searching best crop area", "error", err)
		return errorResponse(filename, err.Error())
	}
	logger.Debug("Best crop", "area", cropArea)

	// Crop image with requested aspect ratio
	type subImager interface {
		SubImage(image.Rectangle) image.Image
	}
	si, ok := img.(subImager)
	if !ok {
		errText := "failed cropping preview subimage"
		logger.Error(errText)
		return errorResponse(filename, errText)
	}
	croppedImage := si.SubImage(cropArea)
	logger.Debug("Cropped image dimensions",
		"width", croppedImage.Bounds().Dx(), "height", croppedImage.Bounds().Dy())

	// Resize image to fit requested params
	resizedImage := resizer.Resize(croppedImage, uint(request.Width), uint(request.Height))

	// In-memory buffer to store JPEG image before we Base64 encode it
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	if err := jpeg.Encode(&buff, resizedImage, &jpeg.Options{Quality: p.quality}); err != nil {
		logger.Error("failed encoding preview to JPEG", "error", err)
		return errorResponse(filename, err.Error())
	}

	// Encode the bytes in the buffer to a base64 string
	preview := make([]byte, base64.StdEncoding.EncodedLen(len(buff.Bytes())))
	base64.StdEncoding.Encode(preview, buff.Bytes())
	logger.Debug("Preview made successfully")

	// Preview made successfully
	return &entities.Response{
		Preview:  preview,
		Filename: filename,
		Status: entities.Status{
			Code: http.StatusOK,
			Text: http.StatusText(http.StatusOK),
		},
	}
}
