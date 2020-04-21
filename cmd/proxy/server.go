package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AcroManiac/micropic/internal/domain/interfaces"

	"github.com/AcroManiac/micropic/internal/domain/entities"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func init() {
	// Turn off debug noise
	gin.SetMode(gin.ReleaseMode)
}

// Server structure
type Server struct {
	host   string
	port   int
	router *gin.Engine
	srv    *http.Server
	cache  interfaces.CacheClient
	rpc    *RMQRPC
}

// NewServer constructs and initializes REST server
func NewServer(host string, port int, cache interfaces.CacheClient, rpc *RMQRPC) *Server {
	server := &Server{
		host:   host,
		port:   port,
		router: gin.Default(),
		srv:    nil,
		cache:  cache,
		rpc:    rpc,
	}

	// Set routing handlers
	server.router.GET("/fill/:width/:height/*imageUrl", server.handlePreview)

	return server
}

func convertString(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("couldn't convert string to int", "error", err)
	}
	return
}

func removeSlash(s string) string {
	if s[0] == '/' {
		return s[1:]
	}
	return s
}

// Get preview from microservices
// Test with:
// curl -ki -X GET -H "Content-Type: image/jpeg" http://localhost:8080/fill/300/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg
func (s *Server) handlePreview(c *gin.Context) {
	// Create params and request
	params := &entities.PreviewParams{
		Width:  convertString(c.Param("width")),
		Height: convertString(c.Param("height")),
		URL:    removeSlash(c.Param("imageUrl")),
	}
	request := &entities.Request{
		Params:  params,
		Headers: c.Request.Header,
	}
	logger.Debug("Incoming request", "request", request)

	// Check preview cache for image first
	preview, err := s.cache.Get(params)
	if err != nil {
		logger.Error("cache request failed", "error", err)

		logger.Debug("Proxying request to previewer")
		// Create RMQRPC context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call previewer and wait for response
		response, err := s.rpc.SendRequest(ctx, request)
		if err != nil {
			logger.Error("previewer RMQRPC request failed", "error", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if response.Status.Code != http.StatusOK {
			logger.Error("error returned from HTTP source", "error", response.Status)
			c.String(response.Status.Code, response.Status.Text)
			return
		}

		// Get preview from previewer response
		preview = response.Preview
	}

	// Decode preview from Base64
	buffSize := base64.StdEncoding.DecodedLen(len(preview.Image))
	decodedPreview := make([]byte, buffSize)
	_, err = base64.StdEncoding.Decode(decodedPreview, preview.Image)
	if err != nil {
		logger.Error("error decoding Base64 to preview", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Extract filename
	filename := request.Params.URL[strings.LastIndex(request.Params.URL, "/")+1:]

	// Return preview file within HTTP response
	reader := bytes.NewReader(decodedPreview)
	contentLength := int64(len(decodedPreview))
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + filename + `"`,
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

// Start HTTP server
func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.router,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "failed starting HTTP server")
	}

	return nil
}

// Stop HTTP server gracefully
func (s *Server) Stop() error {
	if s.srv == nil {
		return errors.New("server object is not created")
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed shutting down RESTful API server")
	}

	return nil
}
