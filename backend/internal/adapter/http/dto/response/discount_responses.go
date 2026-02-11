package response

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/discount"
)

type PromotionRule struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	RuleType      string         `json:"rule_type"`
	Value         float64        `json:"value"`
	Conditions    map[string]any `json:"conditions"`
	FundingSource string         `json:"funding_source"`
	FundingPct    float64        `json:"funding_pct"`
	MaxUsageCount int            `json:"max_usage_count"`
	CurrentUsage  int            `json:"current_usage"`
	ValidFrom     time.Time      `json:"valid_from"`
	ValidTo       time.Time      `json:"valid_to"`
	Status        string         `json:"status"`
	Version       int            `json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func PromotionRuleFromDomain(r *discount.PromotionRule) PromotionRule {
	return PromotionRule{
		ID: r.ID().String(), Name: r.Name(),
		RuleType: r.RuleType(), Value: r.Value(),
		Conditions: r.Conditions(), FundingSource: r.FundingSource(),
		FundingPct: r.FundingPct(), MaxUsageCount: r.MaxUsageCount(),
		CurrentUsage: r.CurrentUsage(), ValidFrom: r.ValidFrom(),
		ValidTo: r.ValidTo(), Status: string(r.Status()),
		Version: r.Version(), CreatedAt: r.CreatedAt(), UpdatedAt: r.UpdatedAt(),
	}
}

type DiscountApplication struct {
	ID                 string    `json:"id"`
	OrderID            string    `json:"order_id"`
	PromotionRuleID    string    `json:"promotion_rule_id"`
	RuleVersion        int       `json:"rule_version"`
	DiscountCents      int64     `json:"discount_cents"`
	Currency           string    `json:"currency"`
	FundingSource      string    `json:"funding_source"`
	SellerShareCents   int64     `json:"seller_share_cents"`
	PlatformShareCents int64     `json:"platform_share_cents"`
	AppliedAt          time.Time `json:"applied_at"`
}

func DiscountApplicationFromDomain(d *discount.DiscountApplication) DiscountApplication {
	return DiscountApplication{
		ID: d.ID().String(), OrderID: d.OrderID().String(),
		PromotionRuleID: d.PromotionRuleID().String(),
		RuleVersion: d.RuleVersion(),
		DiscountCents: d.DiscountCents().Amount(),
		Currency: d.DiscountCents().Currency(),
		FundingSource: d.FundingSource(),
		SellerShareCents: d.SellerShareCents().Amount(),
		PlatformShareCents: d.PlatformShareCents().Amount(),
		AppliedAt: d.AppliedAt(),
	}
}
