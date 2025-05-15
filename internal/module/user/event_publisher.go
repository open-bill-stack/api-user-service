package user

import (
	"context"
	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishUserDelete(ctx context.Context, id uuid.UUID) error
}
