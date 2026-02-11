package product

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo product.Repository
}

func NewService(repo product.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, sku, name, description, category string, unitCost, unitPrice shared.Money) (*product.Product, error) {
	p, err := product.NewProduct(sku, name, unitCost, unitPrice)
	if err != nil {
		return nil, err
	}
	if description != "" || category != "" {
		_ = p.UpdateDetails(name, description, category)
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter product.ListFilter) (*shared.PagedResult[*product.Product], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name, description, category string, unitPrice *shared.Money) (*product.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.UpdateDetails(name, description, category); err != nil {
		return nil, err
	}
	if unitPrice != nil {
		if err := p.UpdatePrice(*unitPrice); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
