package integration

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	Platform   string
	Status     IntegrationStatus
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, i *Integration) error
	GetByID(ctx context.Context, id uuid.UUID) (*Integration, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Integration], error)
	Update(ctx context.Context, i *Integration) error
	Delete(ctx context.Context, id uuid.UUID) error
}
