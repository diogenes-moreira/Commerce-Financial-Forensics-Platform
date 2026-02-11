package cloudaccount

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, a *CloudAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*CloudAccount, error)
	List(ctx context.Context) ([]*CloudAccount, error)
	Update(ctx context.Context, a *CloudAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
}
