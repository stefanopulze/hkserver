package http_rest

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"hkserver/configs"
	"hkserver/internal/registry"
	"sync"
)

func StartHttpService(ctx context.Context, wg *sync.WaitGroup, log *zap.Logger, config configs.Http) {
	wg.Add(1)
	log.Debug(fmt.Sprintf("Starting http server: %d", config.Port))

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version": configs.Version,
			"devices": registry.Size(),
		})
	})

	log.Info(fmt.Sprintf("Started http server: %d", config.Port))
	go app.Listen(fmt.Sprintf(":%d", config.Port))
	go func() {
		<-ctx.Done()
		log.Debug("Closing http server")
		wg.Done()
	}()
}
