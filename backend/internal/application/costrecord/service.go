package costrecord

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo costrecord.Repository
}

func NewService(repo costrecord.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, cloudAccountID uuid.UUID, service string, amount shared.Money, usageDate time.Time) (*costrecord.CostRecord, error) {
	r, err := costrecord.NewCostRecord(cloudAccountID, service, amount, usageDate)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*costrecord.CostRecord, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter costrecord.ListFilter) (*shared.PagedResult[*costrecord.CostRecord], error) {
	return s.repo.List(ctx, filter)
}
