package request

type CreatePromotionRule struct {
	Name          string         `json:"name" binding:"required"`
	RuleType      string         `json:"rule_type" binding:"required"`
	Value         float64        `json:"value" binding:"required"`
	Conditions    map[string]any `json:"conditions"`
	FundingSource string         `json:"funding_source" binding:"required"`
	FundingPct    float64        `json:"funding_pct"`
	MaxUsageCount int            `json:"max_usage_count"`
	ValidFrom     string         `json:"valid_from" binding:"required"`
	ValidTo       string         `json:"valid_to" binding:"required"`
}

type ListPromotionRules struct {
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type ApplyDiscount struct {
	OrderID      string `json:"order_id" binding:"required"`
	RuleID       string `json:"rule_id" binding:"required"`
	SubtotalCents int64 `json:"subtotal_cents" binding:"required"`
	Currency     string `json:"currency"`
}
