package seller

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	Status     Status
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, s *Seller) error
	GetByID(ctx context.Context, id uuid.UUID) (*Seller, error)
	GetByExternalID(ctx context.Context, externalID string) (*Seller, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*Seller], error)
	Update(ctx context.Context, s *Seller) error
}
