package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/AcroManiac/micropic/internal/domain/entities"
)

// HTTPImageLoader structure
type HTTPImageLoader struct{}

// NewHTTPImageLoader constructor
func NewHTTPImageLoader() *HTTPImageLoader {
	return &HTTPImageLoader{}
}

// Load image from HTTP source with request params
func (l *HTTPImageLoader) Load(request *entities.Request) ([]byte, *entities.Status) {
	// Create HTTP request to image source
	req, err := http.NewRequest("GET", "http://"+request.URL, nil)
	if err != nil {
		errText := "error allocating http request"
		logger.Error(errText, "error", err)
		return nil, &entities.Status{
			Code: http.StatusInternalServerError,
			Text: errText,
		}
	}
	// Proxying HTTP headers to request
	for key, value := range request.Headers {
		req.Header.Add(key, strings.Join(value, " "))
	}

	// Make request to image source
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed getting image from HTTP source",
			"url", request.URL, "error", err)
		return nil, &entities.Status{
			Code: resp.StatusCode,
			Text: http.StatusText(resp.StatusCode),
		}
	}
	defer resp.Body.Close()

	// Read image from response
	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("cannot read response body", "error", err)
		return nil, &entities.Status{
			Code: http.StatusNoContent,
			Text: http.StatusText(http.StatusNoContent),
		}
	}

	// TODO: Check for image format

	// Image loaded successfully
	return image, &entities.Status{
		Code: resp.StatusCode,
		Text: http.StatusText(resp.StatusCode),
	}
}
