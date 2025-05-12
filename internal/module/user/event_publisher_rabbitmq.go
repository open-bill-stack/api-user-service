package user

import (
	"api-user-service/internal/module/user/structure"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitEventPublisher struct {
	ch *amqp.Channel
	q  *amqp.Queue
}

func NewEventPublisher(ch *amqp.Channel, q *amqp.Queue) EventPublisher {

	return &rabbitEventPublisher{ch, q}
}

func (r *rabbitEventPublisher) PublishUserDelete(ctx context.Context, uuid string) error {
	var event = structure.Event{
		Type:     "user_delete",
		UserUUID: uuid,
	}

	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.ch.PublishWithContext(ctx,
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate

		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshal,
		})
}
