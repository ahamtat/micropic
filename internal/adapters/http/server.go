package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AcroManiac/micropic/internal/domain/entities"

	"github.com/AcroManiac/micropic/internal/adapters/broker"

	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func init() {
	// Turn off debug noise
	gin.SetMode(gin.ReleaseMode)
}

// Server structure
type Server struct {
	router *gin.Engine
	srv    *http.Server
	mgr    *broker.Manager
}

// NewServer constructs and initializes REST server
func NewServer(mgr *broker.Manager) *Server {
	server := &Server{
		router: gin.Default(),
		srv:    nil,
		mgr:    mgr,
	}

	// Set routing handlers
	server.router.GET("/fill/:width/:height/:imageUrl", server.handlePreview)

	return server
}

func convertString(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		logger.Error("couldn't convert string to int", "error", err)
	}
	return
}

// Get preview from microservices
// Test with:
// curl -ki -X GET -H "Content-Type: application/json" -H "" http://127.0.0.1:2020/api/v3/gateway/configure/6774f85a-0a5b-4059-9b68-9385ecbdcf8e
func (s *Server) handlePreview(c *gin.Context) {
	image := &entities.Image{
		Width:   convertString(c.Param("width")),
		Height:  convertString(c.Param("height")),
		URL:     c.Param("imageUrl"),
		Headers: nil,
	}
	logger.Debug("Image params from incoming HTTP request", "image", image)

	// Check preview cache for image

	//// Make preview and get file path
	//response, err := s.mgr.DoPreviewerRPC(image)
	//if err != nil {
	//	errorText := "gateway RPC request failed"
	//	logger.Error(errorText, "error", err, "gateway", gatewayID)
	//	c.String(http.StatusBadRequest, errorText)
	//	return
	//}
	//if response == nil {
	//	errorText := "no gateway configuration returned"
	//	logger.Error(errorText, "gateway", gatewayID)
	//	c.String(http.StatusBadRequest, errorText)
	//	return
	//}

	//// Return preview file within HTTP response
	//c.JSON(http.StatusOK, response)
}

// Start RESTful server for all interfaces
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", viper.GetInt("rest.port"))
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "failed starting RESTful API server")
	}

	return nil
}

// Stop RESTful API server gracefully
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
