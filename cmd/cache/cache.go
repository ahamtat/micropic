package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/AcroManiac/micropic/internal/adapters/application"
)

func init() {
	application.Init("../../configs/cache.yml")
}

func main() {
	// Create, initialize and start application objects
	app := &appObjects{}
	app.Init()
	app.Start()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Application started. Press Ctrl+C to exit...")

	// Wait for user or OS interrupt
	<-done

	app.Stop()
	logger.Info("Application exited properly")
}

type appObjects struct {
	//rpc
}

func (app *appObjects) Init() {
	// Create and start RPC object
	//app.rpc =
	//if app.rmq == nil {
	//	logger.Fatal("failed creating gRPC server")
	//}
}

func (app *appObjects) Start() {
	// Start RPC loop
	//go app.rpc.Start()
	//logger.Info("gRPC started successfully", "host", host)
}

func (app *appObjects) Stop() {
	// Stop gRPC gracefully
}
