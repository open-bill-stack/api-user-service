package user

import "go.uber.org/fx"

var Module = fx.Module("account",
	fx.Provide(NewHandler),
	fx.Provide(NewService),
	fx.Provide(NewRepository),
)
