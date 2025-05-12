package user

import "go.uber.org/fx"

var Module = fx.Module("account",
	fx.Provide(NewHttpHandler),
	fx.Provide(NewGrpcHandler),

	fx.Provide(NewService),
	fx.Provide(NewRepository),
)
