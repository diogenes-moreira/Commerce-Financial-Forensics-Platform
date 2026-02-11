package discount

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type RuleStatus string

const (
	RuleStatusActive   RuleStatus = "active"
	RuleStatusExpired  RuleStatus = "expired"
	RuleStatusDisabled RuleStatus = "disabled"
)

type PromotionRule struct {
	id            uuid.UUID
	name          string
	ruleType      string
	value         float64
	conditions    map[string]any
	fundingSource string
	fundingPct    float64
	maxUsageCount int
	currentUsage  int
	validFrom     time.Time
	validTo       time.Time
	status        RuleStatus
	version       int
	createdAt     time.Time
	updatedAt     time.Time
}

func NewPromotionRule(
	name, ruleType string, value float64,
	conditions map[string]any, fundingSource string, fundingPct float64,
	maxUsageCount int, validFrom, validTo time.Time,
) (*PromotionRule, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}
	if conditions == nil {
		conditions = make(map[string]any)
	}

	now := time.Now().UTC()
	return &PromotionRule{
		id:            uuid.New(),
		name:          name,
		ruleType:      ruleType,
		value:         value,
		conditions:    conditions,
		fundingSource: fundingSource,
		fundingPct:    fundingPct,
		maxUsageCount: maxUsageCount,
		currentUsage:  0,
		validFrom:     validFrom,
		validTo:       validTo,
		status:        RuleStatusActive,
		version:       1,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func HydratePromotionRule(
	id uuid.UUID, name, ruleType string, value float64,
	conditions map[string]any, fundingSource string, fundingPct float64,
	maxUsageCount, currentUsage int, validFrom, validTo time.Time,
	status RuleStatus, version int, createdAt, updatedAt time.Time,
) *PromotionRule {
	if conditions == nil {
		conditions = make(map[string]any)
	}
	return &PromotionRule{
		id: id, name: name, ruleType: ruleType, value: value,
		conditions: conditions, fundingSource: fundingSource, fundingPct: fundingPct,
		maxUsageCount: maxUsageCount, currentUsage: currentUsage,
		validFrom: validFrom, validTo: validTo, status: status, version: version,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (r *PromotionRule) ID() uuid.UUID            { return r.id }
func (r *PromotionRule) Name() string              { return r.name }
func (r *PromotionRule) RuleType() string          { return r.ruleType }
func (r *PromotionRule) Value() float64            { return r.value }
func (r *PromotionRule) Conditions() map[string]any { return r.conditions }
func (r *PromotionRule) FundingSource() string     { return r.fundingSource }
func (r *PromotionRule) FundingPct() float64       { return r.fundingPct }
func (r *PromotionRule) MaxUsageCount() int        { return r.maxUsageCount }
func (r *PromotionRule) CurrentUsage() int         { return r.currentUsage }
func (r *PromotionRule) ValidFrom() time.Time      { return r.validFrom }
func (r *PromotionRule) ValidTo() time.Time        { return r.validTo }
func (r *PromotionRule) Status() RuleStatus        { return r.status }
func (r *PromotionRule) Version() int              { return r.version }
func (r *PromotionRule) CreatedAt() time.Time      { return r.createdAt }
func (r *PromotionRule) UpdatedAt() time.Time      { return r.updatedAt }

func (r *PromotionRule) IsValid(now time.Time) bool {
	return r.status == RuleStatusActive &&
		!now.Before(r.validFrom) && !now.After(r.validTo)
}

func (r *PromotionRule) CanApply() bool {
	if r.maxUsageCount > 0 && r.currentUsage >= r.maxUsageCount {
		return false
	}
	return r.status == RuleStatusActive
}

func (r *PromotionRule) CalculateDiscount(subtotal shared.Money) shared.Money {
	switch r.ruleType {
	case "percentage":
		return subtotal.Multiply(r.value / 100)
	case "fixed_amount":
		m, err := shared.NewMoney(int64(r.value), subtotal.Currency())
		if err != nil {
			return shared.ZeroMoney(subtotal.Currency())
		}
		return m
	default:
		return shared.ZeroMoney(subtotal.Currency())
	}
}

func (r *PromotionRule) IncrementUsage() {
	r.currentUsage++
	r.updatedAt = time.Now().UTC()
}

func (r *PromotionRule) Disable() {
	r.status = RuleStatusDisabled
	r.updatedAt = time.Now().UTC()
}
