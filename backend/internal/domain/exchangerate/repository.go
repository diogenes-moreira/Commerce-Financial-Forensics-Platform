package exchangerate

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type ListFilter struct {
	BaseCurrency  string
	QuoteCurrency string
	StartDate     *time.Time
	EndDate       *time.Time
	Pagination    shared.Pagination
}

type Repository interface {
	Create(ctx context.Context, r *ExchangeRate) error
	CreateBatch(ctx context.Context, rates []*ExchangeRate) error
	GetByID(ctx context.Context, id uuid.UUID) (*ExchangeRate, error)
	GetLatest(ctx context.Context, baseCurrency, quoteCurrency string) (*ExchangeRate, error)
	GetByDate(ctx context.Context, baseCurrency, quoteCurrency string, date time.Time) (*ExchangeRate, error)
	List(ctx context.Context, filter ListFilter) (*shared.PagedResult[*ExchangeRate], error)
}
