package homekit

import (
	"context"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/registry"
	"log"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup, logger *zap.Logger, opts configs.Homekit) {
	// bridge
	bridge := accessory.NewBridge(accessory.Info{
		Name: opts.Name,
	})

	accessories := registry.GetAllAccessories()

	// configure the ip transport
	config := hc.Config{
		StoragePath: opts.StoragePath,
		Pin:         opts.Pin,
		Port:        opts.Port,
	}
	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		log.Panic(err)
	}

	wg.Add(1)
	go t.Start()

	<-ctx.Done()
	<-t.Stop()
	logger.Debug("Closing hc service")
	wg.Done()
}
