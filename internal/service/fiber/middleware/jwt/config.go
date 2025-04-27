package jwt

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {

	// Optional. Default: nil.
	Validator func(string) bool

	// Unauthorized defines the response body for unauthorized responses.
	// By default it will return with a 401 Unauthorized and the correct WWW-Auth header
	//
	// Optional. Default: nil
	Unauthorized fiber.Handler
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Validator:    nil,
	Unauthorized: nil,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	return cfg
}
