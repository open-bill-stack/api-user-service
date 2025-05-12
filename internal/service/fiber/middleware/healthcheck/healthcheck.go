package healthcheck

import (
	"api-user-service/internal/service/fiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"

	"go.uber.org/fx"
)

type Result struct {
	fx.Out

	Router router.Router `group:"httpRoutes"`
}
type Params struct {
	fx.In
}

type Middleware struct {
}

func (m *Middleware) Register(app *fiber.App) {
	app.Get("/healthcheck", healthcheck.New(healthcheck.Config{
		LivenessProbe:    func(c *fiber.Ctx) bool { return true },
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
}

func NewMiddleware(p Params) (Result, error) {
	return Result{
		Router: &Middleware{},
	}, nil

}
