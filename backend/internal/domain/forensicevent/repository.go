package forensicevent

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, e *ForensicEvent) error
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*ForensicEvent, error)
	ListByType(ctx context.Context, eventType string, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*ForensicEvent], error)
	ListByDateRange(ctx context.Context, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*ForensicEvent], error)
}
