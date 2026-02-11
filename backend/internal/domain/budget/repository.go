package budget

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, b *Budget) error
	GetByID(ctx context.Context, id uuid.UUID) (*Budget, error)
	List(ctx context.Context) ([]*Budget, error)
	Update(ctx context.Context, b *Budget) error
	CreateAlert(ctx context.Context, a *Alert) error
}
