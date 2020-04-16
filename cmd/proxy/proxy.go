package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AcroManiac/micropic/internal/adapters/http"

	"github.com/spf13/viper"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/AcroManiac/micropic/internal/adapters/application"
)

func init() {
	application.Init("../../configs/proxy.yml")
}

func main() {
	// Make cancel context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Create HTTP proxy server
	viper.GetInt("proxy.port")

	// Create RESTful API server
	proxy := http.NewServer(
		viper.GetString("proxy.host"),
		viper.GetInt("proxy.port"),
		nil) //manager)
	if proxy == nil {
		logger.Fatal("could not initialize HTTP server")
	}

	// Start HTTP proxy server in a separate goroutine
	go func() {
		if err := proxy.Start(); err != nil {
			logger.Fatal("could not start HTTP server", "error", err)
		}
	}()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Application started. Press Ctrl+C to exit...")

	// Wait for interruption events
	select {
	case <-ctx.Done():
		logger.Info("Main context cancelled")
	case <-done:
		logger.Info("User or OS interrupted program")
		cancel()
	}

	// Stop HTTP proxy server
	if err := proxy.Stop(); err != nil {
		logger.Error("failed stopping HTTP server")
	}

	logger.Info("Application exited properly")
}
