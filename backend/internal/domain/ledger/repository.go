package ledger

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type DateRange struct {
	Start time.Time
	End   time.Time
}

type Repository interface {
	Create(ctx context.Context, e *LedgerEntry) error
	CreateBatch(ctx context.Context, entries []*LedgerEntry) error
	ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*LedgerEntry, error)
	ListByAccount(ctx context.Context, code string, dateRange DateRange, pagination shared.Pagination) (*shared.PagedResult[*LedgerEntry], error)
	BalanceByAccount(ctx context.Context, code string, dateRange DateRange) (int64, error)
}
