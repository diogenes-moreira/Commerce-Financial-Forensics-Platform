package drift

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/drift"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Service orchestrates drift detection using existing repositories.
type Service struct {
	orderRepo   order.Repository
	productRepo product.Repository
	sellerRepo  seller.Repository
}

// NewService creates a new drift application service.
func NewService(
	orderRepo order.Repository,
	productRepo product.Repository,
	sellerRepo seller.Repository,
) *Service {
	return &Service{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		sellerRepo:  sellerRepo,
	}
}

// DetectForOrder detects price/cost drift for all items in a single order
// by comparing snapshotted values against the current product catalog.
func (s *Service) DetectForOrder(ctx context.Context, orderID uuid.UUID) ([]drift.DriftResult, error) {
	o, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	products, err := s.loadProductsForOrder(ctx, o)
	if err != nil {
		return nil, err
	}

	sl, err := s.sellerRepo.GetByID(ctx, o.SellerID())
	if err != nil {
		return nil, err
	}

	return drift.DetectOrderDrift(o, products, sl)
}

// DetectForPeriod runs drift detection across all orders in a date range
// and returns the aggregated results.
func (s *Service) DetectForPeriod(ctx context.Context, startDate, endDate time.Time) ([]drift.DriftResult, error) {
	var allResults []drift.DriftResult

	page := 1
	pageSize := shared.MaxPageSize

	for {
		filter := order.ListFilter{
			StartDate:  &startDate,
			EndDate:    &endDate,
			Pagination: shared.NewPagination(page, pageSize),
		}

		result, err := s.orderRepo.List(ctx, filter)
		if err != nil {
			return nil, err
		}

		for _, o := range result.Items {
			products, err := s.loadProductsForOrder(ctx, o)
			if err != nil {
				return nil, err
			}

			sl, err := s.sellerRepo.GetByID(ctx, o.SellerID())
			if err != nil {
				return nil, err
			}

			drifts, err := drift.DetectOrderDrift(o, products, sl)
			if err != nil {
				// Skip orders with no items.
				continue
			}

			allResults = append(allResults, drifts...)
		}

		// Check if we have fetched all pages.
		fetched := int64(page * pageSize)
		if fetched >= result.TotalCount {
			break
		}
		page++
	}

	return allResults, nil
}

// loadProductsForOrder fetches the current catalog products for every item in an order.
func (s *Service) loadProductsForOrder(ctx context.Context, o *order.Order) (map[uuid.UUID]*product.Product, error) {
	products := make(map[uuid.UUID]*product.Product)
	for _, item := range o.Items() {
		if _, exists := products[item.ProductID()]; exists {
			continue
		}
		p, err := s.productRepo.GetByID(ctx, item.ProductID())
		if err != nil {
			// Product may have been removed from catalog; skip it.
			continue
		}
		products[item.ProductID()] = p
	}
	return products, nil
}
