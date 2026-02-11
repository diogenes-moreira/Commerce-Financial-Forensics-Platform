package discount

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	ruleRepo discount.PromotionRuleRepository
	appRepo  discount.DiscountApplicationRepository
}

func NewService(ruleRepo discount.PromotionRuleRepository, appRepo discount.DiscountApplicationRepository) *Service {
	return &Service{ruleRepo: ruleRepo, appRepo: appRepo}
}

func (s *Service) CreateRule(
	ctx context.Context,
	name, ruleType string, value float64,
	conditions map[string]any, fundingSource string, fundingPct float64,
	maxUsageCount int, validFrom, validTo time.Time,
) (*discount.PromotionRule, error) {
	rule, err := discount.NewPromotionRule(
		name, ruleType, value, conditions,
		fundingSource, fundingPct, maxUsageCount,
		validFrom, validTo,
	)
	if err != nil {
		return nil, err
	}
	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) GetRule(ctx context.Context, id uuid.UUID) (*discount.PromotionRule, error) {
	return s.ruleRepo.GetByID(ctx, id)
}

func (s *Service) ListRules(ctx context.Context, filter discount.RuleListFilter) (*shared.PagedResult[*discount.PromotionRule], error) {
	return s.ruleRepo.List(ctx, filter)
}

func (s *Service) DisableRule(ctx context.Context, id uuid.UUID) (*discount.PromotionRule, error) {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	rule.Disable()
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) ApplyDiscount(ctx context.Context, orderID uuid.UUID, ruleID uuid.UUID, subtotal shared.Money) (*discount.DiscountApplication, error) {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if !rule.IsValid(now) {
		return nil, discount.ErrRuleExpired
	}
	if !rule.CanApply() {
		return nil, discount.ErrMaxUsageReached
	}

	discountAmount := rule.CalculateDiscount(subtotal)
	application := discount.NewDiscountApplication(orderID, rule, discountAmount)

	rule.IncrementUsage()
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}
	if err := s.appRepo.Create(ctx, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *Service) ListApplicationsByOrder(ctx context.Context, orderID uuid.UUID) ([]*discount.DiscountApplication, error) {
	return s.appRepo.ListByOrder(ctx, orderID)
}
