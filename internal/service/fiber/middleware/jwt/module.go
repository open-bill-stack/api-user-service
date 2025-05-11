package jwt

import (
	"api-user-service/internal/service/jwt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	JWT *jwt.JWT
}

type Result struct {
	fx.Out

	Middleware *CheckJWTMiddleware
}

func NewService(p Params) (Result, error) {
	checkJWTMiddleware := NewCheckJWTMiddleware(Config{
		Validator: func(s string) (*jwt.CustomClaims, bool) {
			claims, err := p.JWT.Verify(s)
			if err != nil {
				return nil, false
			}
			return claims, true
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
	})

	return Result{
		Middleware: checkJWTMiddleware,
	}, nil
}
