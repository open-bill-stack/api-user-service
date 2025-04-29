package jwt

import (
	"api-user-service/internal/service/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Log    *zap.Logger
	Config *config.Config
}
type Result struct {
	fx.Out

	JWT *JWT
}

func NewService(p Params) (Result, error) {
	jwt, _ := NewJWT(p.Config.JWT.PrivateKey, p.Config.JWT.PublicKey)

	return Result{
		JWT: jwt,
	}, nil
}

var Module = fx.Module(
	"JWTModule",
	fx.Provide(
		NewService,
	),
)
