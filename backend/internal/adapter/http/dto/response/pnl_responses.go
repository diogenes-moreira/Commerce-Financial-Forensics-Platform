package response

import (
	"github.com/diogenes/costforensics/backend/internal/application/pnl"
)

// PLRow is the response DTO for a single P&L report row.
type PLRow struct {
	Period          string  `json:"period"`
	Revenue         int64   `json:"revenue"`
	COGS            int64   `json:"cogs"`
	GrossProfit     int64   `json:"gross_profit"`
	DiscountTotal   int64   `json:"discount_total"`
	CommissionTotal int64   `json:"commission_total"`
	ShippingTotal   int64   `json:"shipping_total"`
	NetProfit       int64   `json:"net_profit"`
	GrossMarginPct  float64 `json:"gross_margin_pct"`
	NetMarginPct    float64 `json:"net_margin_pct"`
	OrderCount      int     `json:"order_count"`
	Currency        string  `json:"currency"`
}

// PLRowFromDomain converts a P&L application row to a response DTO.
func PLRowFromDomain(r pnl.PLRow) PLRow {
	return PLRow{
		Period:          r.Period,
		Revenue:         r.Revenue,
		COGS:            r.COGS,
		GrossProfit:     r.GrossProfit,
		DiscountTotal:   r.DiscountTotal,
		CommissionTotal: r.CommissionTotal,
		ShippingTotal:   r.ShippingTotal,
		NetProfit:       r.NetProfit,
		GrossMarginPct:  r.GrossMarginPct,
		NetMarginPct:    r.NetMarginPct,
		OrderCount:      r.OrderCount,
		Currency:        r.Currency,
	}
}

// CohortRow is the response DTO for a single retention cohort row.
type CohortRow struct {
	CohortPeriod  string  `json:"cohort_period"`
	Period        string  `json:"period"`
	CustomerCount int     `json:"customer_count"`
	Revenue       int64   `json:"revenue"`
	RetentionPct  float64 `json:"retention_pct"`
}

// CohortRowFromDomain converts a cohort application row to a response DTO.
func CohortRowFromDomain(r pnl.CohortRow) CohortRow {
	return CohortRow{
		CohortPeriod:  r.CohortPeriod,
		Period:        r.Period,
		CustomerCount: r.CustomerCount,
		Revenue:       r.Revenue,
		RetentionPct:  r.RetentionPct,
	}
}
