package jwt

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type CheckJWTMiddleware struct {
	cfg *Config
}

func NewCheckJWTMiddleware(config ...Config) *CheckJWTMiddleware {
	cfg := configDefault(config...)
	return &CheckJWTMiddleware{&cfg}
}

func (h *CheckJWTMiddleware) GetHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// Get authorization header
		auth := c.Get(fiber.HeaderAuthorization)

		// Check if the header contains content besides "basic".
		if len(auth) <= 7 || !utils.EqualFold(auth[:7], "bearer ") {
			return h.cfg.Unauthorized(c)
		}

		token := auth[7:]

		h.cfg.Validator(token)
		return c.Next()
	}
}
