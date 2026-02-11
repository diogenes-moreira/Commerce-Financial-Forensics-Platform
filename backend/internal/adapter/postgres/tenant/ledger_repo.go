package tenant

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LedgerRepo struct {
	db *gorm.DB
}

func NewLedgerRepo(db *gorm.DB) *LedgerRepo {
	return &LedgerRepo{db: db}
}

func (r *LedgerRepo) Create(ctx context.Context, e *ledger.LedgerEntry) error {
	return r.db.WithContext(ctx).Create(mapper.LedgerEntryToModel(e)).Error
}

func (r *LedgerRepo) CreateBatch(ctx context.Context, entries []*ledger.LedgerEntry) error {
	if len(entries) == 0 {
		return nil
	}
	models := make([]model.LedgerEntry, len(entries))
	for i, e := range entries {
		models[i] = *mapper.LedgerEntryToModel(e)
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

func (r *LedgerRepo) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*ledger.LedgerEntry, error) {
	var models []model.LedgerEntry
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).
		Order("effective_date ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]*ledger.LedgerEntry, len(models))
	for i := range models {
		items[i] = mapper.LedgerEntryToDomain(&models[i])
	}
	return items, nil
}

func (r *LedgerRepo) ListByAccount(ctx context.Context, code string, dateRange ledger.DateRange, pagination shared.Pagination) (*shared.PagedResult[*ledger.LedgerEntry], error) {
	q := r.db.WithContext(ctx).Model(&model.LedgerEntry{}).Where("account_code = ?", code)

	if !dateRange.Start.IsZero() {
		q = q.Where("effective_date >= ?", dateRange.Start)
	}
	if !dateRange.End.IsZero() {
		q = q.Where("effective_date <= ?", dateRange.End)
	}

	var total int64
	q.Count(&total)

	var models []model.LedgerEntry
	if err := q.Offset(pagination.Offset()).Limit(pagination.PageSize()).
		Order("effective_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*ledger.LedgerEntry, len(models))
	for i := range models {
		items[i] = mapper.LedgerEntryToDomain(&models[i])
	}

	return &shared.PagedResult[*ledger.LedgerEntry]{
		Items:      items,
		TotalCount: total,
		Page:       pagination.Page(),
		PageSize:   pagination.PageSize(),
	}, nil
}

func (r *LedgerRepo) BalanceByAccount(ctx context.Context, code string, dateRange ledger.DateRange) (int64, error) {
	q := r.db.WithContext(ctx).Model(&model.LedgerEntry{}).Where("account_code = ?", code)

	if !dateRange.Start.IsZero() {
		q = q.Where("effective_date >= ?", dateRange.Start)
	}
	if !dateRange.End.IsZero() {
		q = q.Where("effective_date <= ?", dateRange.End)
	}

	var result struct {
		Balance int64
	}
	err := q.Select("COALESCE(SUM(CASE WHEN side = 'debit' THEN amount_cents ELSE -amount_cents END), 0) as balance").
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Balance, nil
}
