package router

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Params struct {
	fx.In

	App    *grpc.Server
	Log    *zap.Logger
	Router []Router `group:"routes"`
}

func RunRoute(lc fx.Lifecycle, p Params) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, r := range p.Router {
				r.Register(p.App)
			}
			p.Log.Debug("Include routes")
			return nil
		},
	})

}
