package application

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

// HealthChecker structure.
type HealthChecker struct {
	router  *gin.Engine
	isReady *atomic.Value
	port    int
	srv     *http.Server
}

// NewHealthChecker constructor.
func NewHealthChecker(router *gin.Engine) *HealthChecker {
	if router == nil {
		router = gin.Default()
	}
	isReady := &atomic.Value{}
	isReady.Store(false)

	// Create health checker object
	checker := &HealthChecker{
		router:  router,
		isReady: isReady,
		port:    0,
		srv:     nil,
	}

	// Add kubernetes endpoints
	router.GET("/healthz", checker.handleHealth)
	router.GET("/readyz", checker.handleReady)
	return checker
}

// SetReady sets application ready flag
func (hc *HealthChecker) SetReady() {
	hc.isReady.Store(true)
}

// handleHealth is a liveness probe.
func (hc HealthChecker) handleHealth(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// readyz is a readiness probe.
func (hc *HealthChecker) handleReady(ctx *gin.Context) {
	if hc.isReady == nil || !hc.isReady.Load().(bool) {
		ctx.String(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
		return
	}
	ctx.Status(http.StatusOK)
}

// HealthCheckerServer structure.
type HealthCheckerServer struct {
	Chk  *HealthChecker
	port int
	srv  *http.Server
}

// NewHealthCheckerServer constructor.
func NewHealthCheckerServer(port int) *HealthCheckerServer {
	// Turn off debug noise
	gin.SetMode(gin.ReleaseMode)

	chk := NewHealthChecker(nil)
	chk.port = port
	return &HealthCheckerServer{
		Chk:  chk,
		port: port,
		srv:  nil,
	}
}

// Start HTTP server.
func (s *HealthCheckerServer) Start() error {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.Chk.router,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "failed starting HTTP server")
	}

	return nil
}

// Stop HTTP server gracefully.
func (s *HealthCheckerServer) Stop() error {
	if s.srv == nil {
		return errors.New("server object was not created")
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed shutting down HTTP server")
	}

	return nil
}
