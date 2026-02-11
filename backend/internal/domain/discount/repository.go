package discount

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type RuleListFilter struct {
	Status     RuleStatus
	Pagination shared.Pagination
}

type PromotionRuleRepository interface {
	Create(ctx context.Context, r *PromotionRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*PromotionRule, error)
	List(ctx context.Context, filter RuleListFilter) (*shared.PagedResult[*PromotionRule], error)
	Update(ctx context.Context, r *PromotionRule) error
}

type DiscountApplicationRepository interface {
	Create(ctx context.Context, a *DiscountApplication) error
	ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*DiscountApplication, error)
}
