package cloudaccount

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/google/uuid"
)

type Service struct {
	repo cloudaccount.Repository
}

func NewService(repo cloudaccount.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, provider cloudaccount.Provider, name, externalID string) (*cloudaccount.CloudAccount, error) {
	a, err := cloudaccount.NewCloudAccount(provider, name, externalID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*cloudaccount.CloudAccount, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*cloudaccount.CloudAccount, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name string) (*cloudaccount.CloudAccount, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.Rename(name); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) BeginSync(ctx context.Context, id uuid.UUID) (*cloudaccount.CloudAccount, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.BeginSync(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}
