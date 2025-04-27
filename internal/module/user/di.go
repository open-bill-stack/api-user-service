package user

import (
	httpRouter "api-user-service/internal/service/fiber/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Log     *zap.Logger
	Service *Service
}

type Result struct {
	fx.Out

	Router httpRouter.Router `group:"routes"`
}
