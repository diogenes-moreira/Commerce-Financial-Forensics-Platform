package pnl

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/margin"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Granularity determines how orders are bucketed into time periods.
type Granularity string

const (
	GranularityDay     Granularity = "day"
	GranularityMonth   Granularity = "month"
	GranularityQuarter Granularity = "quarter"
	GranularityYear    Granularity = "year"
)

// ParseGranularity converts a string to a Granularity, defaulting to month.
func ParseGranularity(s string) Granularity {
	switch Granularity(s) {
	case GranularityDay, GranularityMonth, GranularityQuarter, GranularityYear:
		return Granularity(s)
	default:
		return GranularityMonth
	}
}

// PLRow represents a single row in a P&L report for a specific time period.
type PLRow struct {
	Period          string
	Revenue         int64
	COGS            int64
	GrossProfit     int64
	DiscountTotal   int64
	CommissionTotal int64
	ShippingTotal   int64
	NetProfit       int64
	GrossMarginPct  float64
	NetMarginPct    float64
	OrderCount      int
	Currency        string
}

// CohortRow represents a single row in a retention cohort report.
type CohortRow struct {
	CohortPeriod  string
	Period        string
	CustomerCount int
	Revenue       int64
	RetentionPct  float64
}

// Service orchestrates P&L and cohort report generation.
type Service struct {
	orderRepo    order.Repository
	sellerRepo   seller.Repository
	discountRepo discount.DiscountApplicationRepository
}

// NewService creates a new P&L application service.
func NewService(
	orderRepo order.Repository,
	sellerRepo seller.Repository,
	discountRepo discount.DiscountApplicationRepository,
) *Service {
	return &Service{
		orderRepo:    orderRepo,
		sellerRepo:   sellerRepo,
		discountRepo: discountRepo,
	}
}

// GeneratePL aggregates orders by time period and produces P&L rows with
// margin calculations for each bucket.
func (s *Service) GeneratePL(
	ctx context.Context,
	startDate, endDate time.Time,
	granularity Granularity,
	currency string,
) ([]PLRow, error) {
	if currency == "" {
		currency = "USD"
	}

	orders, err := s.fetchAllOrders(ctx, startDate, endDate, nil)
	if err != nil {
		return nil, err
	}

	// Cache sellers to avoid redundant fetches.
	sellerCache := make(map[uuid.UUID]*seller.Seller)

	// Accumulate P&L data into period buckets.
	buckets := make(map[string]*PLRow)

	for _, o := range orders {
		if o.Currency() != currency {
			continue
		}

		period := formatPeriod(o.OrderDate(), granularity)

		sl, err := s.getCachedSeller(ctx, o.SellerID(), sellerCache)
		if err != nil {
			return nil, err
		}

		discounts, err := s.discountRepo.ListByOrder(ctx, o.ID())
		if err != nil {
			return nil, err
		}

		breakdown, err := margin.CalculateMargin(o, sl, discounts)
		if err != nil {
			return nil, err
		}

		row, ok := buckets[period]
		if !ok {
			row = &PLRow{
				Period:   period,
				Currency: currency,
			}
			buckets[period] = row
		}

		row.Revenue += breakdown.RevenueCents
		row.COGS += breakdown.COGSCents
		row.DiscountTotal += breakdown.DiscountCents
		row.CommissionTotal += breakdown.CommissionCents
		row.ShippingTotal += breakdown.ShippingCostCents
		row.OrderCount++
	}

	// Compute derived fields and collect into a sorted slice.
	result := make([]PLRow, 0, len(buckets))
	for _, row := range buckets {
		row.GrossProfit = row.Revenue - row.COGS
		row.NetProfit = row.Revenue - row.COGS - row.CommissionTotal - row.ShippingTotal

		if row.Revenue > 0 {
			row.GrossMarginPct = float64(row.GrossProfit) / float64(row.Revenue) * 100
			row.NetMarginPct = float64(row.NetProfit) / float64(row.Revenue) * 100
		}

		result = append(result, *row)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Period < result[j].Period
	})

	return result, nil
}

// GenerateCohorts groups customers by the period of their first order and tracks
// their ordering behavior in subsequent periods.
func (s *Service) GenerateCohorts(
	ctx context.Context,
	startDate, endDate time.Time,
	granularity Granularity,
) ([]CohortRow, error) {
	orders, err := s.fetchAllOrders(ctx, startDate, endDate, nil)
	if err != nil {
		return nil, err
	}

	// Determine each customer's first order period (their cohort).
	// customerFirstOrder maps customerID -> earliest order date.
	customerFirstOrder := make(map[uuid.UUID]time.Time)
	for _, o := range orders {
		existing, ok := customerFirstOrder[o.CustomerID()]
		if !ok || o.OrderDate().Before(existing) {
			customerFirstOrder[o.CustomerID()] = o.OrderDate()
		}
	}

	// Map each customer to their cohort period string.
	customerCohort := make(map[uuid.UUID]string)
	cohortSizes := make(map[string]int)
	for customerID, firstDate := range customerFirstOrder {
		cohort := formatPeriod(firstDate, granularity)
		customerCohort[customerID] = cohort
		cohortSizes[cohort]++
	}

	// For each order, record the customer's activity in the order's period.
	// Key: "cohortPeriod|period", value: set of customerIDs and revenue.
	type bucketKey struct {
		cohort string
		period string
	}
	type bucketData struct {
		customers map[uuid.UUID]struct{}
		revenue   int64
	}
	buckets := make(map[bucketKey]*bucketData)

	for _, o := range orders {
		customerID := o.CustomerID()
		cohort := customerCohort[customerID]
		period := formatPeriod(o.OrderDate(), granularity)

		key := bucketKey{cohort: cohort, period: period}
		data, ok := buckets[key]
		if !ok {
			data = &bucketData{customers: make(map[uuid.UUID]struct{})}
			buckets[key] = data
		}
		data.customers[customerID] = struct{}{}
		data.revenue += o.TotalCents().Amount()
	}

	// Build result rows.
	result := make([]CohortRow, 0, len(buckets))
	for key, data := range buckets {
		customerCount := len(data.customers)
		originalSize := cohortSizes[key.cohort]

		var retentionPct float64
		if originalSize > 0 {
			retentionPct = float64(customerCount) / float64(originalSize) * 100
		}

		result = append(result, CohortRow{
			CohortPeriod:  key.cohort,
			Period:        key.period,
			CustomerCount: customerCount,
			Revenue:       data.revenue,
			RetentionPct:  retentionPct,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CohortPeriod != result[j].CohortPeriod {
			return result[i].CohortPeriod < result[j].CohortPeriod
		}
		return result[i].Period < result[j].Period
	})

	return result, nil
}

// fetchAllOrders paginates through the order repository and returns all orders
// in the specified date range.
func (s *Service) fetchAllOrders(
	ctx context.Context,
	startDate, endDate time.Time,
	sellerID *uuid.UUID,
) ([]*order.Order, error) {
	var allOrders []*order.Order
	page := 1
	pageSize := shared.MaxPageSize

	for {
		filter := order.ListFilter{
			StartDate:  &startDate,
			EndDate:    &endDate,
			SellerID:   sellerID,
			Pagination: shared.NewPagination(page, pageSize),
		}

		result, err := s.orderRepo.List(ctx, filter)
		if err != nil {
			return nil, err
		}

		allOrders = append(allOrders, result.Items...)

		fetched := int64(page * pageSize)
		if fetched >= result.TotalCount {
			break
		}
		page++
	}

	return allOrders, nil
}

// getCachedSeller retrieves a seller from the cache or fetches from the repo.
func (s *Service) getCachedSeller(
	ctx context.Context,
	sellerID uuid.UUID,
	cache map[uuid.UUID]*seller.Seller,
) (*seller.Seller, error) {
	if sl, ok := cache[sellerID]; ok {
		return sl, nil
	}
	sl, err := s.sellerRepo.GetByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	cache[sellerID] = sl
	return sl, nil
}

// formatPeriod formats a time into a period string based on granularity.
func formatPeriod(t time.Time, g Granularity) string {
	switch g {
	case GranularityDay:
		return t.Format("2006-01-02")
	case GranularityMonth:
		return t.Format("2006-01")
	case GranularityQuarter:
		q := (t.Month()-1)/3 + 1
		return fmt.Sprintf("%d-Q%d", t.Year(), q)
	case GranularityYear:
		return t.Format("2006")
	default:
		return t.Format("2006-01")
	}
}
