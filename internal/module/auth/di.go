package auth

import (
	"api-user-service/internal/module/user"
	jwtMiddleware "api-user-service/internal/service/fiber/middleware/jwt"
	httpRouter "api-user-service/internal/service/fiber/router"
	"api-user-service/internal/service/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Log                *zap.Logger
	JWT                *jwt.JWT
	CheckJWTMiddleware *jwtMiddleware.CheckJWTMiddleware
	UserService        *user.Service
}

type Result struct {
	fx.Out

	Router httpRouter.Router `group:"httpRoutes"`
}
