package costreport

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, r *CostReport) error
	GetByID(ctx context.Context, id uuid.UUID) (*CostReport, error)
	List(ctx context.Context) ([]*CostReport, error)
}
