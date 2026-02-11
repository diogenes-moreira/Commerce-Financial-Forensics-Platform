package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(ctx context.Context, p *payment.Payment) error {
	return r.db.WithContext(ctx).Create(mapper.PaymentToModel(p)).Error
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	var m model.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payment.ErrPaymentNotFound
		}
		return nil, err
	}
	return mapper.PaymentToDomain(&m), nil
}

func (r *PaymentRepo) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*payment.Payment, error) {
	var models []model.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).
		Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]*payment.Payment, len(models))
	for i := range models {
		items[i] = mapper.PaymentToDomain(&models[i])
	}
	return items, nil
}

func (r *PaymentRepo) List(ctx context.Context, filter payment.ListFilter) (*shared.PagedResult[*payment.Payment], error) {
	q := r.db.WithContext(ctx).Model(&model.Payment{})

	if filter.OrderID != nil {
		q = q.Where("order_id = ?", *filter.OrderID)
	}
	if filter.Direction != "" {
		q = q.Where("direction = ?", string(filter.Direction))
	}
	if filter.CounterpartyID != nil {
		q = q.Where("counterparty_id = ?", *filter.CounterpartyID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}

	var total int64
	q.Count(&total)

	var models []model.Payment
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*payment.Payment, len(models))
	for i := range models {
		items[i] = mapper.PaymentToDomain(&models[i])
	}

	return &shared.PagedResult[*payment.Payment]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *PaymentRepo) Update(ctx context.Context, p *payment.Payment) error {
	return r.db.WithContext(ctx).Save(mapper.PaymentToModel(p)).Error
}
