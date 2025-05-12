package database

import (
	"api-user-service/internal/service/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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
	DatabaseConn *pgxpool.Pool
}

func NewDatabase(p Params) (Result, error) {
	pool, _ := pgxpool.New(context.Background(), p.Config.Database.Url)
	return Result{
		DatabaseConn: pool,
	}, nil
}

var Module = fx.Module(
	"DatabasePgxModule",
	fx.Provide(
		NewDatabase,
	),
)
