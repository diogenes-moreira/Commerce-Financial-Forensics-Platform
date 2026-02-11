package tenant

import (
	"context"
	"errors"
	"time"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/exchangerate"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExchangeRateRepo struct {
	db *gorm.DB
}

func NewExchangeRateRepo(db *gorm.DB) *ExchangeRateRepo {
	return &ExchangeRateRepo{db: db}
}

func (r *ExchangeRateRepo) Create(ctx context.Context, rate *exchangerate.ExchangeRate) error {
	return r.db.WithContext(ctx).Create(mapper.ExchangeRateToModel(rate)).Error
}

func (r *ExchangeRateRepo) CreateBatch(ctx context.Context, rates []*exchangerate.ExchangeRate) error {
	models := make([]model.ExchangeRate, len(rates))
	for i, rate := range rates {
		models[i] = *mapper.ExchangeRateToModel(rate)
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

func (r *ExchangeRateRepo) GetByID(ctx context.Context, id uuid.UUID) (*exchangerate.ExchangeRate, error) {
	var m model.ExchangeRate
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exchangerate.ErrRateNotFound
		}
		return nil, err
	}
	return mapper.ExchangeRateToDomain(&m), nil
}

func (r *ExchangeRateRepo) GetLatest(ctx context.Context, baseCurrency, quoteCurrency string) (*exchangerate.ExchangeRate, error) {
	var m model.ExchangeRate
	if err := r.db.WithContext(ctx).
		Where("base_currency = ? AND quote_currency = ?", baseCurrency, quoteCurrency).
		Order("effective_date DESC").First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exchangerate.ErrRateNotFound
		}
		return nil, err
	}
	return mapper.ExchangeRateToDomain(&m), nil
}

func (r *ExchangeRateRepo) GetByDate(ctx context.Context, baseCurrency, quoteCurrency string, date time.Time) (*exchangerate.ExchangeRate, error) {
	var m model.ExchangeRate
	if err := r.db.WithContext(ctx).
		Where("base_currency = ? AND quote_currency = ? AND effective_date <= ?", baseCurrency, quoteCurrency, date).
		Order("effective_date DESC").First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exchangerate.ErrRateNotFound
		}
		return nil, err
	}
	return mapper.ExchangeRateToDomain(&m), nil
}

func (r *ExchangeRateRepo) List(ctx context.Context, filter exchangerate.ListFilter) (*shared.PagedResult[*exchangerate.ExchangeRate], error) {
	q := r.db.WithContext(ctx).Model(&model.ExchangeRate{})

	if filter.BaseCurrency != "" {
		q = q.Where("base_currency = ?", filter.BaseCurrency)
	}
	if filter.QuoteCurrency != "" {
		q = q.Where("quote_currency = ?", filter.QuoteCurrency)
	}
	if filter.StartDate != nil {
		q = q.Where("effective_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		q = q.Where("effective_date <= ?", *filter.EndDate)
	}

	var total int64
	q.Count(&total)

	var models []model.ExchangeRate
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("effective_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*exchangerate.ExchangeRate, len(models))
	for i := range models {
		items[i] = mapper.ExchangeRateToDomain(&models[i])
	}

	return &shared.PagedResult[*exchangerate.ExchangeRate]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}
