package costrecord

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	CloudAccountID *uuid.UUID
	Service        string
	Category       Category
	StartDate      *time.Time
	EndDate        *time.Time
	Pagination     shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, r *CostRecord) error
	GetByID(ctx context.Context, id uuid.UUID) (*CostRecord, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*CostRecord], error)
}
