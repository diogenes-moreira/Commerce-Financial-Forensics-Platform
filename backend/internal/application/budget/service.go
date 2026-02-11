package budget

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/budget"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo budget.Repository
}

func NewService(repo budget.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string, amount shared.Money, periodStart, periodEnd time.Time, alertThresholdPct int) (*budget.Budget, error) {
	b, err := budget.NewBudget(name, amount, periodStart, periodEnd, alertThresholdPct)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*budget.Budget, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name string) (*budget.Budget, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := b.Rename(name); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) RecordSpend(ctx context.Context, id uuid.UUID, amount shared.Money) (*budget.Budget, *budget.Alert, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	alert, err := b.RecordSpend(amount)
	if err != nil {
		return nil, nil, err
	}
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, nil, err
	}
	if alert != nil {
		if err := s.repo.CreateAlert(ctx, alert); err != nil {
			return b, nil, err
		}
	}
	return b, alert, nil
}
