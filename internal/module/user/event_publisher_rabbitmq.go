package user

import (
	"api-user-service/internal/module/user/structure"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitEventPublisher struct {
	ch *amqp.Channel
}

func NewEventPublisher(ch *amqp.Channel) EventPublisher {
	return &rabbitEventPublisher{ch}
}

func (r *rabbitEventPublisher) PublishUserDelete(ctx context.Context, id uuid.UUID) error {
	var event = structure.Event{
		UserUUID: id.String(),
	}

	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.ch.PublishWithContext(ctx,
		"user.events",  // exchange
		"user.deleted", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshal,
		})
}
