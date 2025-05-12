package user

import (
	httpRouter "api-user-service/internal/service/fiber/router"
	grpcRouter "api-user-service/internal/service/grpc/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Log     *zap.Logger
	Service *Service
}

type HttpResult struct {
	fx.Out

	Router httpRouter.Router `group:"httpRoutes"`
}

type GrpcResult struct {
	fx.Out

	Router grpcRouter.Router `group:"grpcRoutes"`
}
