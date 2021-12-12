package discover

import (
	"context"
	"go.uber.org/zap"
	"hkserver/internal/discover/shellies"
	"hkserver/internal/transport"
)

type Service interface {
	Init(ctx context.Context) error
	Scan(ctx context.Context) error
}

func Scan(ctx context.Context, logger *zap.Logger) {
	logger.Info("Starting discovery devices on network")

	shs := shellies.New(transport.Mqtt, logger)
	shs.Init(ctx)
	shs.Scan(ctx)

}
