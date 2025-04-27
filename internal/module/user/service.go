package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetCredentialsByEmail(ctx context.Context, email string) (*Credential, error) {
	return s.repo.GetCredentialsByEmail(ctx, email)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) Create(ctx context.Context, input CreateUser) (*User, error) {
	return s.repo.Create(ctx, input)
}

func (s *Service) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email)
}
