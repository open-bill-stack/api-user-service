package middleware

import (
	"api-user-service/internal/service/fiber/middleware/healthcheck"
	"api-user-service/internal/service/fiber/middleware/jwt"
	middlewareRecover "api-user-service/internal/service/fiber/middleware/recover"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"MiddlewareModule",
	fx.Provide(
		healthcheck.NewMiddleware,
		middlewareRecover.NewMiddleware,
		jwt.NewService,
	),
)
