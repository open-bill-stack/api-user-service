package user

import (
	pb "api-user-service/internal/module/user/grpc"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandle struct {
	pb.UnimplementedUserServiceServer

	log     *zap.Logger
	service *Service
}

func NewGrpcHandler(p Params) (GrpcResult, error) {
	return GrpcResult{
		Router: &GrpcHandle{
			log:     p.Log,
			service: p.Service,
		},
	}, nil
}

func (h *GrpcHandle) Register(app *grpc.Server) {
	pb.RegisterUserServiceServer(app, h)
}

func (h *GrpcHandle) ExistsByID(ctx context.Context, data *pb.ExistsByIDRequest) (*pb.ExistsByIDResponse, error) {
	exists, err := h.service.ExistsByUUID(ctx, data.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking uuid: %v", err)
	}
	return &pb.ExistsByIDResponse{Exists: exists}, nil
}
