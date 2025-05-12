package user

import (
	"context"
)

type EventPublisher interface {
	PublishUserDelete(ctx context.Context, uuid string) error
}
