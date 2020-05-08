package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AcroManiac/micropic/internal/adapters/file"
	"github.com/AcroManiac/micropic/internal/domain/interfaces"
	"github.com/AcroManiac/micropic/internal/domain/usecases"

	"github.com/AcroManiac/micropic/internal/adapters/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/spf13/viper"

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
	cache      interfaces.Cache
	lsnr       net.Listener
	grpcServer *grpc.Server
	hltServer  *application.HealthCheckerServer
}

func (app *appObjects) Init() {
	// Create cache object
	app.cache = usecases.NewLRUCache(
		viper.GetInt("cache.size"),
		file.NewFileStorage(
			viper.GetString("cache.dirname")))

	// Create a gRPC Server with gRPC interceptor
	app.grpcServer = grpc.NewServer()
	grpcapi.RegisterCacheServer(app.grpcServer, grpcapi.NewCacheServerImpl(app.cache))

	// Register reflection service on gRPC server
	reflection.Register(app.grpcServer)

	// Create health checking server
	app.hltServer = application.NewHealthCheckerServer(viper.GetInt("health.port"))
}

func (app *appObjects) Start() {
	// Create listener for gRPC server
	var err error
	app.lsnr, err = net.Listen("tcp", fmt.Sprintf("%s:%d",
		viper.GetString("grpc.ip"),
		viper.GetInt("grpc.port")))
	if err != nil {
		logger.Fatal("failed to listen tcp", "error", err)
	}

	// Listen gRPC server
	logger.Info("Starting gRPC server...")
	go func() {
		if err := app.grpcServer.Serve(app.lsnr); err != nil {
			logger.Fatal("error while starting gRPC server", "error", err)
		}
	}()

	// Start health checking server
	logger.Info("Starting HealthCheck server...", "port", viper.GetInt("health.port"))
	go func() {
		if err := app.hltServer.Start(); err != nil {
			logger.Error("error starting HealthCheck server", "error", err)
			return
		}
	}()
	app.hltServer.Chk.SetReady()
}

func (app *appObjects) Stop() {
	// Make gRPC server graceful shutdown
	app.grpcServer.GracefulStop()

	// Clear file cache
	_ = app.cache.Clean()

	// Stop health checking server
	if err := app.hltServer.Stop(); err != nil {
		logger.Error("error stopping HealthCheck server", "error", err)
	}
}
