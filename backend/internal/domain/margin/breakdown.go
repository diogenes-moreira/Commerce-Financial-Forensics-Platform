package margin

import (
	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// MarginBreakdown is a computed result representing the full margin analysis
// for a single order. It is NOT a stored entity.
type MarginBreakdown struct {
	OrderID               uuid.UUID
	RevenueCents          int64
	COGSCents             int64
	DiscountCents         int64
	SellerDiscountCents   int64
	PlatformDiscountCents int64
	CommissionCents       int64
	ShippingCostCents     int64
	TaxCents              int64
	GrossMarginCents      int64
	NetMarginCents        int64
	RealMarginPct         float64
	Currency              string
}

// CalculateMargin is a pure domain service that computes a full margin breakdown
// from an order, its seller, and any discount applications.
func CalculateMargin(o *order.Order, s *seller.Seller, discounts []*discount.DiscountApplication) (MarginBreakdown, error) {
	if o == nil {
		return MarginBreakdown{}, ErrOrderRequired
	}

	currency := o.Currency()

	// Compute total COGS from order items.
	totalCOGS := shared.ZeroMoney(currency)
	for _, item := range o.Items() {
		totalCOGS, _ = totalCOGS.Add(item.COGSCents())
	}

	// Compute commission from seller.
	var commissionMoney shared.Money
	if s != nil {
		commissionMoney = s.CalculateCommission(o.SubtotalCents())
	} else {
		commissionMoney = shared.ZeroMoney(currency)
	}

	// Aggregate discount splits from discount applications.
	totalDiscount := shared.ZeroMoney(currency)
	sellerDiscount := shared.ZeroMoney(currency)
	platformDiscount := shared.ZeroMoney(currency)
	for _, d := range discounts {
		totalDiscount, _ = totalDiscount.Add(d.DiscountCents())
		sellerDiscount, _ = sellerDiscount.Add(d.SellerShareCents())
		platformDiscount, _ = platformDiscount.Add(d.PlatformShareCents())
	}

	revenue := o.TotalCents().Amount()
	cogs := totalCOGS.Amount()
	commission := commissionMoney.Amount()
	shipping := o.ShippingCents().Amount()

	grossMargin := revenue - cogs
	netMargin := revenue - cogs - commission - shipping

	var realMarginPct float64
	if revenue > 0 {
		realMarginPct = float64(netMargin) / float64(revenue) * 100
	}

	return MarginBreakdown{
		OrderID:               o.ID(),
		RevenueCents:          revenue,
		COGSCents:             cogs,
		DiscountCents:         totalDiscount.Amount(),
		SellerDiscountCents:   sellerDiscount.Amount(),
		PlatformDiscountCents: platformDiscount.Amount(),
		CommissionCents:       commission,
		ShippingCostCents:     shipping,
		TaxCents:              o.TaxCents().Amount(),
		GrossMarginCents:      grossMargin,
		NetMarginCents:        netMargin,
		RealMarginPct:         realMarginPct,
		Currency:              currency,
	}, nil
}
