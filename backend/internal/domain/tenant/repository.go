package tenant

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, t *Tenant) error
	Delete(ctx context.Context, id uuid.UUID) error
}
