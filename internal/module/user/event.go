package user

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ParamsRun struct {
	fx.In

	Log         *zap.Logger
	AQMPChannel *amqp.Channel
}

func RunEvent(lc fx.Lifecycle, p ParamsRun) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return p.AQMPChannel.ExchangeDeclare(
				"user.events", // name
				"topic",
				true,
				false,
				false,
				false,
				nil, // arguments
			)
		},
	})
}
