package order

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo order.Repository
}

func NewService(repo order.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, externalID string, sellerID, customerID uuid.UUID, currency string) (*order.Order, error) {
	o, err := order.NewOrder(externalID, sellerID, customerID, currency)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter order.ListFilter) (*shared.PagedResult[*order.Order], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) AddItem(ctx context.Context, orderID, productID uuid.UUID, sku, name string, qty int, unitPrice, unitCost shared.Money) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if err := o.AddItem(productID, sku, name, qty, unitPrice, unitCost); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Confirm(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Confirm(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Ship(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Ship(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Deliver(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Deliver(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Cancel(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Cancel(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Refund(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Refund(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}
