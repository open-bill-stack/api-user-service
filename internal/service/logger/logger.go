package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"LoggerModule",
	fx.Provide(
		zap.NewProduction,
	),
)
