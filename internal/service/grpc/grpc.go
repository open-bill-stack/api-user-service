package grpc

import (
	"api-user-service/internal/service/config"
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Params struct {
	fx.In

	GrpcApp      *grpc.Server
	GrpcListener *net.Listener
	Log          *zap.Logger
	Config       *config.Config
}
type Result struct {
	fx.Out
	GrpcApp      *grpc.Server
	GrpcListener *net.Listener
}

func NewGrpcApp() (Result, error) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":50051"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return Result{
		GrpcApp:      grpcServer,
		GrpcListener: &grpcListener,
	}, nil
}

func RunGrpcApp(lc fx.Lifecycle, p Params) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				p.Log.Debug("starting gRPC server")
				if err := (*p.GrpcApp).Serve(*p.GrpcListener); err != nil {
					p.Log.Panic("Error starting GRPC server:", zap.Error(err))
				}
			}()
			p.Log.Info(fmt.Sprintf("GRPC server started on:%d", 50051))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Graceful shutdown Fiber сервера
			p.Log.Info("Shutting down GRPC server...")
			return (*p.GrpcListener).Close()
		},
	})

}

var Module = fx.Module(
	"GrpcAppModule",
	fx.Provide(
		NewGrpcApp,
	),
	fx.Invoke(
		RunGrpcApp,
	),
)
