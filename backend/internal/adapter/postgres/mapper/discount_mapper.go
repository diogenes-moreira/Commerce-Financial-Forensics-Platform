package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
)

// PromotionRule

func PromotionRuleToModel(r *discount.PromotionRule) *model.PromotionRule {
	condJSON, _ := json.Marshal(r.Conditions())
	return &model.PromotionRule{
		ID:            r.ID(),
		Name:          r.Name(),
		RuleType:      r.RuleType(),
		Value:         r.Value(),
		Conditions:    condJSON,
		FundingSource: r.FundingSource(),
		FundingPct:    r.FundingPct(),
		MaxUsageCount: r.MaxUsageCount(),
		CurrentUsage:  r.CurrentUsage(),
		ValidFrom:     r.ValidFrom(),
		ValidTo:       r.ValidTo(),
		Status:        string(r.Status()),
		Version:       r.Version(),
		CreatedAt:     r.CreatedAt(),
		UpdatedAt:     r.UpdatedAt(),
	}
}

func PromotionRuleToDomain(m *model.PromotionRule) *discount.PromotionRule {
	var conditions map[string]any
	_ = json.Unmarshal(m.Conditions, &conditions)
	if conditions == nil {
		conditions = make(map[string]any)
	}
	return discount.HydratePromotionRule(
		m.ID, m.Name, m.RuleType, m.Value,
		conditions, m.FundingSource, m.FundingPct,
		m.MaxUsageCount, m.CurrentUsage, m.ValidFrom, m.ValidTo,
		discount.RuleStatus(m.Status), m.Version, m.CreatedAt, m.UpdatedAt,
	)
}

// DiscountApplication

func DiscountApplicationToModel(d *discount.DiscountApplication) *model.DiscountApplication {
	return &model.DiscountApplication{
		ID:                 d.ID(),
		OrderID:            d.OrderID(),
		PromotionRuleID:    d.PromotionRuleID(),
		RuleVersion:        d.RuleVersion(),
		DiscountCents:      d.DiscountCents().Amount(),
		Currency:           d.DiscountCents().Currency(),
		FundingSource:      d.FundingSource(),
		SellerShareCents:   d.SellerShareCents().Amount(),
		PlatformShareCents: d.PlatformShareCents().Amount(),
		AppliedAt:          d.AppliedAt(),
	}
}

func DiscountApplicationToDomain(m *model.DiscountApplication) *discount.DiscountApplication {
	discountAmt := shared.MustNewMoney(m.DiscountCents, m.Currency)
	sellerShare := shared.MustNewMoney(m.SellerShareCents, m.Currency)
	platformShare := shared.MustNewMoney(m.PlatformShareCents, m.Currency)
	return discount.HydrateDiscountApplication(
		m.ID, m.OrderID, m.PromotionRuleID, m.RuleVersion,
		discountAmt, sellerShare, platformShare,
		m.FundingSource, m.AppliedAt,
	)
}
