package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, input CreateUser) (*User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*Credential, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUUID(ctx context.Context, uuid string) (bool, error)
	DeleteByUUID(ctx context.Context, uuid string) (bool, error)
}
