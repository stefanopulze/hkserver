package transport

import (
	"context"
	"fmt"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/transport/mqtt"
	"sync"
)

var Mqtt pahoMqtt.Client

func Initialize(ctx context.Context, wg *sync.WaitGroup, config *configs.Config, logger *zap.Logger) {
	logger.Debug("Initializing transports")

	mc, err := mqtt.New(ctx, config.Mqtt, wg)
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot initialize mqtt client: %s", err.Error()))
	} else {
		Mqtt = mc
	}

	logger.Info("Transports initialized")
}
