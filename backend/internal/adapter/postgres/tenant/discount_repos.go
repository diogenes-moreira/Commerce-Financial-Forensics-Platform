package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PromotionRuleRepo

type PromotionRuleRepo struct {
	db *gorm.DB
}

func NewPromotionRuleRepo(db *gorm.DB) *PromotionRuleRepo {
	return &PromotionRuleRepo{db: db}
}

func (r *PromotionRuleRepo) Create(ctx context.Context, rule *discount.PromotionRule) error {
	return r.db.WithContext(ctx).Create(mapper.PromotionRuleToModel(rule)).Error
}

func (r *PromotionRuleRepo) GetByID(ctx context.Context, id uuid.UUID) (*discount.PromotionRule, error) {
	var m model.PromotionRule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, discount.ErrRuleNotFound
		}
		return nil, err
	}
	return mapper.PromotionRuleToDomain(&m), nil
}

func (r *PromotionRuleRepo) List(ctx context.Context, filter discount.RuleListFilter) (*shared.PagedResult[*discount.PromotionRule], error) {
	q := r.db.WithContext(ctx).Model(&model.PromotionRule{})

	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}

	var total int64
	q.Count(&total)

	var models []model.PromotionRule
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*discount.PromotionRule, len(models))
	for i := range models {
		items[i] = mapper.PromotionRuleToDomain(&models[i])
	}

	return &shared.PagedResult[*discount.PromotionRule]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *PromotionRuleRepo) Update(ctx context.Context, rule *discount.PromotionRule) error {
	return r.db.WithContext(ctx).Save(mapper.PromotionRuleToModel(rule)).Error
}

// DiscountApplicationRepo

type DiscountApplicationRepo struct {
	db *gorm.DB
}

func NewDiscountApplicationRepo(db *gorm.DB) *DiscountApplicationRepo {
	return &DiscountApplicationRepo{db: db}
}

func (r *DiscountApplicationRepo) Create(ctx context.Context, a *discount.DiscountApplication) error {
	return r.db.WithContext(ctx).Create(mapper.DiscountApplicationToModel(a)).Error
}

func (r *DiscountApplicationRepo) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*discount.DiscountApplication, error) {
	var models []model.DiscountApplication
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).
		Order("applied_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]*discount.DiscountApplication, len(models))
	for i := range models {
		items[i] = mapper.DiscountApplicationToDomain(&models[i])
	}
	return items, nil
}
