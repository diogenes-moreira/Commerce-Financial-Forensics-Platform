package response

import (
	"github.com/diogenes/costforensics/backend/internal/domain/drift"
	"github.com/diogenes/costforensics/backend/internal/domain/margin"
)

// MarginBreakdown is the response DTO for a computed margin analysis.
type MarginBreakdown struct {
	OrderID               string  `json:"order_id"`
	RevenueCents          int64   `json:"revenue_cents"`
	COGSCents             int64   `json:"cogs_cents"`
	DiscountCents         int64   `json:"discount_cents"`
	SellerDiscountCents   int64   `json:"seller_discount_cents"`
	PlatformDiscountCents int64   `json:"platform_discount_cents"`
	CommissionCents       int64   `json:"commission_cents"`
	ShippingCostCents     int64   `json:"shipping_cost_cents"`
	TaxCents              int64   `json:"tax_cents"`
	GrossMarginCents      int64   `json:"gross_margin_cents"`
	NetMarginCents        int64   `json:"net_margin_cents"`
	RealMarginPct         float64 `json:"real_margin_pct"`
	Currency              string  `json:"currency"`
}

// MarginBreakdownFromDomain converts a domain MarginBreakdown to a response DTO.
func MarginBreakdownFromDomain(m *margin.MarginBreakdown) MarginBreakdown {
	return MarginBreakdown{
		OrderID:               m.OrderID.String(),
		RevenueCents:          m.RevenueCents,
		COGSCents:             m.COGSCents,
		DiscountCents:         m.DiscountCents,
		SellerDiscountCents:   m.SellerDiscountCents,
		PlatformDiscountCents: m.PlatformDiscountCents,
		CommissionCents:       m.CommissionCents,
		ShippingCostCents:     m.ShippingCostCents,
		TaxCents:              m.TaxCents,
		GrossMarginCents:      m.GrossMarginCents,
		NetMarginCents:        m.NetMarginCents,
		RealMarginPct:         m.RealMarginPct,
		Currency:              m.Currency,
	}
}

// DriftResult is the response DTO for a price/cost drift detection result.
type DriftResult struct {
	OrderID       string  `json:"order_id"`
	ItemID        string  `json:"item_id"`
	FieldName     string  `json:"field_name"`
	ExpectedValue int64   `json:"expected_value"`
	ActualValue   int64   `json:"actual_value"`
	DriftPct      float64 `json:"drift_pct"`
	Severity      string  `json:"severity"`
	Currency      string  `json:"currency"`
}

// DriftResultFromDomain converts a domain DriftResult to a response DTO.
func DriftResultFromDomain(d drift.DriftResult) DriftResult {
	return DriftResult{
		OrderID:       d.OrderID.String(),
		ItemID:        d.ItemID.String(),
		FieldName:     d.FieldName,
		ExpectedValue: d.ExpectedValue,
		ActualValue:   d.ActualValue,
		DriftPct:      d.DriftPct,
		Severity:      d.Severity,
		Currency:      d.Currency,
	}
}
