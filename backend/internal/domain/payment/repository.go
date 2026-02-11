package payment

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	OrderID        *uuid.UUID
	Direction      Direction
	CounterpartyID *uuid.UUID
	Status         Status
	Pagination     shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*Payment, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Payment], error)
	Update(ctx context.Context, p *Payment) error
}
