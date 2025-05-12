package fiber

import (
	"api-user-service/internal/service/config"
	"api-user-service/internal/service/fiber/router"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	HttpApp *fiber.App
	Log     *zap.Logger
	Config  *config.Config
	Router  []router.Router `group:"httpRoutes"`
}
type Result struct {
	fx.Out
	HttpApp *fiber.App
}

func NewFiberApp() (Result, error) {
	configFiber := fiber.Config{
		StrictRouting:         true,
		DisableStartupMessage: false,
		EnablePrintRoutes:     true,
	}
	return Result{
		HttpApp: fiber.New(configFiber),
	}, nil
}

func RunFiberApp(lc fx.Lifecycle, p Params) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, r := range p.Router {
				r.Register(p.HttpApp)
			}
			go func() {
				if err := p.HttpApp.Listen(fmt.Sprintf(":%d", p.Config.App.HttpPort)); err != nil {
					p.Log.Panic("Error starting Fiber server:", zap.Error(err))
				}
			}()
			p.Log.Info(fmt.Sprintf("Fiber server started on:%d", p.Config.App.HttpPort))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Graceful shutdown Fiber сервера
			p.Log.Info("Shutting down Fiber server...")
			return p.HttpApp.Shutdown()
		},
	})

}

var Module = fx.Module(
	"FiberAppModule",
	fx.Provide(
		NewFiberApp,
	),
	fx.Invoke(
		RunFiberApp,
	),
)
