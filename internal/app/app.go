package app

import (
	"api-user-service/internal/module/auth"
	"api-user-service/internal/module/user"
	"api-user-service/internal/service/config"
	"api-user-service/internal/service/database"
	"api-user-service/internal/service/fiber"
	"api-user-service/internal/service/fiber/middleware"
	"api-user-service/internal/service/grpc"
	"api-user-service/internal/service/jwt"
	"api-user-service/internal/service/logger"
	"api-user-service/internal/service/rabbitmq"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func Run(cmd *cobra.Command) {
	fx.New(
		fx.Provide(func() *cobra.Command { return cmd }),

		logger.Module,
		config.Module,
		database.Module,
		jwt.Module,
		rabbitmq.Module,

		// global middleware
		middleware.Module,

		// fiber http
		user.Module,
		auth.Module,

		// fiber
		fiber.Module,

		// grpc
		grpc.Module,
	).Run()
}
