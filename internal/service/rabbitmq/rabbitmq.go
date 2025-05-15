package rabbitmq

import (
	"api-user-service/internal/service/config"
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ParamsAMQP struct {
	fx.In

	Log    *zap.Logger
	Config *config.Config
}
type ParamsChannel struct {
	fx.In

	Log      *zap.Logger
	Config   *config.Config
	AQMPConn *amqp.Connection
}

type ParamsAMQPRun struct {
	fx.In

	Log      *zap.Logger
	Config   *config.Config
	AQMPConn *amqp.Connection
}
type ParamsChannelRun struct {
	fx.In

	Log         *zap.Logger
	Config      *config.Config
	AQMPChannel *amqp.Channel
}
type ResultAMQP struct {
	fx.Out
	AQMPConn *amqp.Connection
}
type ResultChannel struct {
	fx.Out
	AQMPChannel *amqp.Channel
}

func NewAMQP(p ParamsAMQP) (ResultAMQP, error) {
	conn, err := amqp.Dial(p.Config.AMQP.Url)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	return ResultAMQP{
		AQMPConn: conn,
	}, nil
}
func NewChannel(p ParamsChannel) (ResultChannel, error) {
	ch, err := p.AQMPConn.Channel()
	if err != nil {
		log.Panicf("Failed to load channel to RabbitMQ: %s", err)
	}

	return ResultChannel{
		AQMPChannel: ch,
	}, nil
}

func RunAMQPConnection(lc fx.Lifecycle, p ParamsAMQPRun) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Graceful shutdown Fiber сервера
			p.Log.Info("Shutting down AMQP client...")
			return p.AQMPConn.Close()
		},
	})
}

func RunAMQPChannel(lc fx.Lifecycle, p ParamsChannelRun) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Graceful shutdown Fiber сервера
			p.Log.Info("Shutting down AMQP channel...")
			return p.AQMPChannel.Close()
		},
	})
}

var Module = fx.Module(
	"AMQPModule",
	fx.Provide(
		NewAMQP,
		NewChannel,
	),
	fx.Invoke(
		RunAMQPConnection,
		RunAMQPChannel,
	),
)
