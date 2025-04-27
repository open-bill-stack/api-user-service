package recover

import (
	"api-user-service/internal/service/fiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out

	Router router.Router `group:"routes"`
}

type Middleware struct{}

func (*Middleware) Register(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
}

func NewMiddleware() (Result, error) {
	return Result{
		Router: &Middleware{},
	}, nil

}
