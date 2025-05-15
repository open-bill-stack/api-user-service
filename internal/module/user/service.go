package user

import (
	"context"
	"github.com/google/uuid"
)

type Service struct {
	repo           Repository
	eventPublisher EventPublisher
}

func NewService(r Repository, e EventPublisher) *Service {
	return &Service{
		repo:           r,
		eventPublisher: e,
	}
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
func (s *Service) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}
func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) (bool, error) {
	if status, err := s.repo.DeleteByID(ctx, id); err != nil {
		return status, err
	}
	if err := s.eventPublisher.PublishUserDelete(ctx, id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}
