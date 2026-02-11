package anomaly

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/anomaly"
	"github.com/google/uuid"
)

type Service struct {
	repo anomaly.Repository
}

func NewService(repo anomaly.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*anomaly.CostAnomaly, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*anomaly.CostAnomaly, error) {
	return s.repo.List(ctx)
}

func (s *Service) Resolve(ctx context.Context, id uuid.UUID) (*anomaly.CostAnomaly, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.Resolve(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}
