package customer

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo customer.Repository
}

func NewService(repo customer.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, externalID, name, email, segment string) (*customer.Customer, error) {
	c, err := customer.NewCustomer(externalID, name, email, segment)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter customer.ListFilter) (*shared.PagedResult[*customer.Customer], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name, email, segment string) (*customer.Customer, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := c.UpdateDetails(name, email, segment); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
