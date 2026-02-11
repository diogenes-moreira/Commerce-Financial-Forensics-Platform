package drift

import (
	"math"

	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/google/uuid"
)

// DriftResult is a computed comparison between snapshotted values in order items
// and current catalog/seller data. It is NOT a stored entity.
type DriftResult struct {
	OrderID       uuid.UUID
	ItemID        uuid.UUID
	FieldName     string  // "unit_price", "unit_cost", "commission_pct"
	ExpectedValue int64   // current catalog value
	ActualValue   int64   // snapshotted value in order
	DriftPct      float64 // (actual - expected) / expected * 100
	Severity      string  // low | medium | high | critical
	Currency      string
}

// classifySeverity determines drift severity from the absolute drift percentage.
func classifySeverity(absDriftPct float64) string {
	switch {
	case absDriftPct > 30:
		return "critical"
	case absDriftPct > 15:
		return "high"
	case absDriftPct > 5:
		return "medium"
	default:
		return "low"
	}
}

// DetectOrderDrift is a pure domain function that compares snapshotted order item
// values against the current product catalog and seller commission data.
func DetectOrderDrift(o *order.Order, products map[uuid.UUID]*product.Product, s *seller.Seller) ([]DriftResult, error) {
	if o == nil || len(o.Items()) == 0 {
		return nil, ErrNoOrderItems
	}

	var results []DriftResult
	currency := o.Currency()

	for _, item := range o.Items() {
		p, ok := products[item.ProductID()]
		if !ok {
			// Product not found in catalog; skip drift analysis for this item.
			continue
		}

		// Check unit_price drift.
		catalogPrice := p.UnitPriceCents().Amount()
		snappedPrice := item.UnitPriceCents().Amount()
		if catalogPrice != snappedPrice && catalogPrice > 0 {
			driftPct := float64(snappedPrice-catalogPrice) / float64(catalogPrice) * 100
			absDrift := math.Abs(driftPct)
			results = append(results, DriftResult{
				OrderID:       o.ID(),
				ItemID:        item.ID(),
				FieldName:     "unit_price",
				ExpectedValue: catalogPrice,
				ActualValue:   snappedPrice,
				DriftPct:      driftPct,
				Severity:      classifySeverity(absDrift),
				Currency:      currency,
			})
		}

		// Check unit_cost drift.
		catalogCost := p.UnitCostCents().Amount()
		snappedCost := item.UnitCostCents().Amount()
		if catalogCost != snappedCost && catalogCost > 0 {
			driftPct := float64(snappedCost-catalogCost) / float64(catalogCost) * 100
			absDrift := math.Abs(driftPct)
			results = append(results, DriftResult{
				OrderID:       o.ID(),
				ItemID:        item.ID(),
				FieldName:     "unit_cost",
				ExpectedValue: catalogCost,
				ActualValue:   snappedCost,
				DriftPct:      driftPct,
				Severity:      classifySeverity(absDrift),
				Currency:      currency,
			})
		}
	}

	return results, nil
}
