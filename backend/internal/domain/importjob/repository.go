package importjob

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	Status     ImportStatus
	EntityType string
	Pagination shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, j *ImportJob) error
	GetByID(ctx context.Context, id uuid.UUID) (*ImportJob, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*ImportJob], error)
	Update(ctx context.Context, j *ImportJob) error
}
