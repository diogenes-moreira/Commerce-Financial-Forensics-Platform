package customer

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	Segment    string
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, c *Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	GetByExternalID(ctx context.Context, externalID string) (*Customer, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Customer], error)
	Update(ctx context.Context, c *Customer) error
}
