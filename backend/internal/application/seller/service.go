package seller

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo seller.Repository
}

func NewService(repo seller.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, externalID, code, name, email string, commissionPct float64) (*seller.Seller, error) {
	sl, err := seller.NewSeller(externalID, code, name, email, commissionPct)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, sl); err != nil {
		return nil, err
	}
	return sl, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*seller.Seller, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter seller.ListFilter) (*shared.PagedResult[*seller.Seller], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name, email string, commissionPct *float64) (*seller.Seller, error) {
	sl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := sl.UpdateDetails(name, email); err != nil {
		return nil, err
	}
	if commissionPct != nil {
		if err := sl.UpdateCommission(*commissionPct); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Update(ctx, sl); err != nil {
		return nil, err
	}
	return sl, nil
}
