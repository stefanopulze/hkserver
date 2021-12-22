package main

import (
	"context"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/di"
	"hkserver/internal/discover"
	"hkserver/internal/homekit"
	"hkserver/internal/http_rest"
	"hkserver/internal/logger"
	"hkserver/internal/registry"
	"hkserver/internal/transport"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	config, err := configs.Initialize()
	if err != nil {
		panic(err)
	}

	log := logger.InitializeMainLogger(config.LogLevel)
	log.Info("Starting hkserver")

	// main context
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// handle sig terms
	go handleShutdownSigs(cancel, log)

	// service
	di.Cron = gocron.NewScheduler(time.UTC)

	// initialize http server
	http_rest.StartHttpService(ctx, &wg, log, config.Http)

	// initialize transports
	transport.Initialize(ctx, &wg, config, log)

	// initialize device registry
	registry.Initialize(
		ctx,
		log.With(zap.String("component", "registry")),
		config.Devices,
	)

	// Discovery service
	if config.Discovery.Enable {
		discover.Scan(ctx, log)
	}

	// TODO implement discovery refresh
	if config.Homekit.Enable {
		go func() {
			if config.Discovery.Enable {
				time.Sleep(5 * time.Second)
			}
			homekit.Start(ctx, &wg, log, config.Homekit)
		}()
	}

	wg.Wait()
	log.Info("Shutdown")
}

// handleSigs capture exit event fo graceful shutdown
func handleShutdownSigs(cancel context.CancelFunc, log *zap.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Debug("Captured signal", zap.String("sig", sig.String()))
	cancel()
}
