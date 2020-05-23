package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahamtat/micropic/internal/adapters/grpcapi"

	"github.com/ahamtat/micropic/internal/domain/interfaces"

	"github.com/ahamtat/micropic/internal/adapters/broker"

	"github.com/spf13/viper"

	"github.com/ahamtat/micropic/internal/adapters/logger"

	"github.com/ahamtat/micropic/internal/adapters/application"
)

func init() {
	application.Init("../../configs/proxy.yml")
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
				// TODO: Restart RMQRPC
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
	cache   interfaces.CacheClient
	rpc     *RMQRPC
	proxy   *Server
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

	// Create and start RMQRPC object
	app.rpc = NewRPC(app.manager.Conn)
	if app.rpc == nil {
		logger.Fatal("failed creating RabbitMQ RMQRPC object")
	}

	// Create cache client
	app.cache = grpcapi.NewCacheClientImpl(
		viper.GetString("grpc.host"),
		viper.GetInt("grpc.port"))

	// Create HTTP proxy server
	app.proxy = NewServer(
		viper.GetString("proxy.host"),
		viper.GetInt("proxy.port"),
		viper.GetInt("proxy.timeout"),
		app.cache,
		app.rpc)
	if app.proxy == nil {
		logger.Fatal("could not initialize HTTP server")
	}

	// Create health checking server for Kubernetes
	app.hServer = application.NewHealthCheckerServer(viper.GetInt("health.port"))
}

func (app *appObjects) Start() {
	// Start RMQRPC loop
	go app.rpc.Start()
	logger.Info("RabbitMQ RMQRPC object started successfully")

	// Start broker connection listener
	go app.manager.ConnectionListener(app.ctx)

	// Start HTTP proxy server in a separate goroutine
	go func() {
		if err := app.proxy.Start(); err != nil {
			logger.Fatal("could not start HTTP server", "error", err)
		}
	}()

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
	// Stop HTTP proxy server
	if err := app.proxy.Stop(); err != nil {
		logger.Error("failed stopping HTTP server")
	}

	// Stop RabbitMQ connection
	app.cancel()
	app.rpc.Stop()
	if err := app.manager.Close(); err != nil {
		logger.Error("failed stopping RabbitMQ connection", "error", err)
	}

	// Stop health checking server
	if err := app.hServer.Stop(); err != nil {
		logger.Error("error stopping HealthCheck server", "error", err)
	}
}
