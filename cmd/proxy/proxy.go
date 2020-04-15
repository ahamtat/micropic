package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	logger.Info("Application exited properly")
}
