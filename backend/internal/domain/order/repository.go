package order

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	SellerID   *uuid.UUID
	CustomerID *uuid.UUID
	Status     Status
	StartDate  *time.Time
	EndDate    *time.Time
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	GetByExternalID(ctx context.Context, externalID string) (*Order, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Order], error)
	Update(ctx context.Context, o *Order) error
}
