package grpc

import (
	"api-user-service/internal/service/config"
	"api-user-service/internal/service/grpc/router"
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strings"
)

type ParamsServer struct {
	fx.In

	Log    *zap.Logger
	Config *config.Config
}
type ParamsServerRun struct {
	fx.In

	GrpcApp      *grpc.Server
	GrpcListener *net.Listener
	Log          *zap.Logger
	Config       *config.Config
	Router       []router.Router `group:"grpcRoutes"`
}
type ResultServer struct {
	fx.Out
	GrpcApp      *grpc.Server
	GrpcListener *net.Listener
}

func NewGrpcServerApp(p ParamsServer) (ResultServer, error) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.Config.App.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return ResultServer{
		GrpcApp:      grpcServer,
		GrpcListener: &grpcListener,
	}, nil
}

func RunGrpcServerApp(lc fx.Lifecycle, p ParamsServerRun) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, r := range p.Router {
				r.Register(p.GrpcApp)
			}
			go func() {
				p.Log.Debug("starting gRPC server")
				if err := (*p.GrpcApp).Serve(*p.GrpcListener); err != nil && !errors.Is(err, net.ErrClosed) {
					if strings.Contains(err.Error(), "use of closed network connection") {
						p.Log.Error("GRPC server stopped with error", zap.Error(err))
					}
				}
			}()
			p.Log.Info(fmt.Sprintf("GRPC server started on:%d", p.Config.App.GrpcPort))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Graceful shutdown Fiber сервера
			p.Log.Info("Shutting down GRPC server...")
			return (*p.GrpcListener).Close()
		},
	})

}

var ModuleServer = fx.Module(
	"GrpcServerAppModule",
	fx.Provide(
		NewGrpcServerApp,
	),
	fx.Invoke(
		RunGrpcServerApp,
	),
)
