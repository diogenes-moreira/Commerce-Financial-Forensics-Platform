package anomaly

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a *CostAnomaly) error
	GetByID(ctx context.Context, id uuid.UUID) (*CostAnomaly, error)
	List(ctx context.Context) ([]*CostAnomaly, error)
	Update(ctx context.Context, a *CostAnomaly) error
}
