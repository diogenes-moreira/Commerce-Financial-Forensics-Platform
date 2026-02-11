package product

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	Category   string
	Status     Status
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Product], error)
	Update(ctx context.Context, p *Product) error
}
