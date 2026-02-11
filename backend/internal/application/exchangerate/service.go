package exchangerate

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo exchangerate.Repository
}

func NewService(repo exchangerate.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, baseCurrency, quoteCurrency string, rate float64, source string, effectiveDate time.Time) (*exchangerate.ExchangeRate, error) {
	r, err := exchangerate.NewExchangeRate(baseCurrency, quoteCurrency, rate, source, effectiveDate)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*exchangerate.ExchangeRate, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetLatest(ctx context.Context, baseCurrency, quoteCurrency string) (*exchangerate.ExchangeRate, error) {
	return s.repo.GetLatest(ctx, baseCurrency, quoteCurrency)
}

func (s *Service) List(ctx context.Context, filter exchangerate.ListFilter) (*shared.PagedResult[*exchangerate.ExchangeRate], error) {
	return s.repo.List(ctx, filter)
}
