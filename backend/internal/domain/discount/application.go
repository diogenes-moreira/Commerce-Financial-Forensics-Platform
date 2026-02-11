package discount

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type DiscountApplication struct {
	id                 uuid.UUID
	orderID            uuid.UUID
	promotionRuleID    uuid.UUID
	ruleVersion        int
	discountCents      shared.Money
	fundingSource      string
	sellerShareCents   shared.Money
	platformShareCents shared.Money
	appliedAt          time.Time
}

func NewDiscountApplication(
	orderID uuid.UUID, rule *PromotionRule, discountAmount shared.Money,
) *DiscountApplication {
	var sellerShare, platformShare shared.Money
	currency := discountAmount.Currency()

	switch rule.FundingSource() {
	case "seller":
		sellerShare = discountAmount
		platformShare = shared.ZeroMoney(currency)
	case "platform":
		sellerShare = shared.ZeroMoney(currency)
		platformShare = discountAmount
	case "shared":
		sellerShare = discountAmount.Multiply(rule.FundingPct() / 100)
		platformShare = discountAmount.Multiply((100 - rule.FundingPct()) / 100)
	default:
		sellerShare = shared.ZeroMoney(currency)
		platformShare = discountAmount
	}

	return &DiscountApplication{
		id:                 uuid.New(),
		orderID:            orderID,
		promotionRuleID:    rule.ID(),
		ruleVersion:        rule.Version(),
		discountCents:      discountAmount,
		fundingSource:      rule.FundingSource(),
		sellerShareCents:   sellerShare,
		platformShareCents: platformShare,
		appliedAt:          time.Now().UTC(),
	}
}

func HydrateDiscountApplication(
	id, orderID, promotionRuleID uuid.UUID, ruleVersion int,
	discount, sellerShare, platformShare shared.Money,
	fundingSource string, appliedAt time.Time,
) *DiscountApplication {
	return &DiscountApplication{
		id: id, orderID: orderID, promotionRuleID: promotionRuleID,
		ruleVersion: ruleVersion, discountCents: discount,
		fundingSource: fundingSource, sellerShareCents: sellerShare,
		platformShareCents: platformShare, appliedAt: appliedAt,
	}
}

func (d *DiscountApplication) ID() uuid.UUID               { return d.id }
func (d *DiscountApplication) OrderID() uuid.UUID           { return d.orderID }
func (d *DiscountApplication) PromotionRuleID() uuid.UUID   { return d.promotionRuleID }
func (d *DiscountApplication) RuleVersion() int             { return d.ruleVersion }
func (d *DiscountApplication) DiscountCents() shared.Money  { return d.discountCents }
func (d *DiscountApplication) FundingSource() string        { return d.fundingSource }
func (d *DiscountApplication) SellerShareCents() shared.Money   { return d.sellerShareCents }
func (d *DiscountApplication) PlatformShareCents() shared.Money { return d.platformShareCents }
func (d *DiscountApplication) AppliedAt() time.Time         { return d.appliedAt }
