package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahamtat/micropic/internal/adapters/broker"

	"github.com/spf13/viper"

	"github.com/ahamtat/micropic/internal/adapters/logger"

	"github.com/ahamtat/micropic/internal/adapters/application"
)

func init() {
	application.Init("../../configs/previewer.yml")
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

OUTER:
	for {
		select {
		// Wait for user or OS interrupt
		case <-done:
			break OUTER

		// Catch broker connection notification
		case connErr := <-app.manager.Done:
			if connErr != nil {
				// Call context to stop i/o operations and scheduler
				app.cancel()

				// Recreate broker connection and scheduler
				if err := app.manager.Reconnect(app.ctx); err != nil {
					logger.Error("error reconnecting RabbitMQ", "error", err)
					break OUTER
				}
				// TODO: Restart RMQ Listener
			}
		}
	}

	app.Stop()
	logger.Info("Application exited properly")
}

type appObjects struct {
	ctx     context.Context
	cancel  context.CancelFunc
	manager *broker.Manager
	rmq     *RMQListener
	hServer *application.HealthCheckerServer
}

func (app *appObjects) Init() {
	// Make cancel context
	app.ctx, app.cancel = context.WithCancel(context.Background())

	// Create broker manager
	app.manager = broker.NewManager(
		viper.GetString("amqp.protocol"),
		viper.GetString("amqp.user"),
		viper.GetString("amqp.password"),
		viper.GetString("amqp.host"),
		viper.GetInt("amqp.port"))
	if app.manager == nil {
		logger.Fatal("failed connecting to RabbitMQ")
		//return // to prevent linter warning
	}
	logger.Info("RabbitMQ broker connected", "host", viper.GetString("amqp.host"))

	// Create and start RPC object
	app.rmq = NewRMQListener(app.manager.Conn,
		viper.GetInt("previewer.quality"),
		viper.GetString("grpc.host"),
		viper.GetInt("grpc.port"))
	if app.rmq == nil {
		logger.Fatal("failed creating RabbitMQ server")
	}

	// Create health checking server for Kubernetes
	app.hServer = application.NewHealthCheckerServer(viper.GetInt("health.port"))
}

func (app *appObjects) Start() {
	// Start RPC loop
	go app.rmq.Start()
	logger.Info("RabbitMQ listener started successfully")

	// Start broker connection listener
	go app.manager.ConnectionListener(app.ctx)

	// Start health checking server
	logger.Info("Starting HealthCheck server...", "port", viper.GetInt("health.port"))
	go func() {
		if err := app.hServer.Start(); err != nil {
			logger.Error("error starting HealthCheck server", "error", err)
			return
		}
	}()
	app.hServer.Chk.SetReady()
}

func (app *appObjects) Stop() {
	// Stop RabbitMQ connection
	app.cancel()
	app.rmq.Stop()
	if err := app.manager.Close(); err != nil {
		logger.Error("failed stopping RabbitMQ connection", "error", err)
	}

	// Stop health checking server
	if err := app.hServer.Stop(); err != nil {
		logger.Error("error stopping HealthCheck server", "error", err)
	}
}
