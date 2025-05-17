package cors

import (
	"api-user-service/internal/service/fiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out

	Router router.Router `group:"httpRoutes"`
}

type Middleware struct{}

func (*Middleware) Register(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
}

func NewMiddleware() (Result, error) {
	return Result{
		Router: &Middleware{},
	}, nil

}
