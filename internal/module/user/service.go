package user

import (
	"context"
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
func (s *Service) ExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	return s.repo.ExistsByUUID(ctx, uuid)
}
func (s *Service) DeleteByUUID(ctx context.Context, uuid string) (bool, error) {
	if status, err := s.repo.DeleteByUUID(ctx, uuid); err != nil {
		return status, err
	}
	if err := s.eventPublisher.PublishUserDelete(ctx, uuid); err != nil {
		return false, err
	}
	return true, nil
}
