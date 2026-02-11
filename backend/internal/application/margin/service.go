package margin

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/discount"
	"github.com/diogenes/costforensics/backend/internal/domain/margin"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Service orchestrates margin calculation using existing repositories.
type Service struct {
	orderRepo    order.Repository
	sellerRepo   seller.Repository
	discountRepo discount.DiscountApplicationRepository
}

// NewService creates a new margin application service.
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

// CalculateForOrder computes the full margin breakdown for a single order.
func (s *Service) CalculateForOrder(ctx context.Context, orderID uuid.UUID) (*margin.MarginBreakdown, error) {
	o, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	sl, err := s.sellerRepo.GetByID(ctx, o.SellerID())
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

	return &breakdown, nil
}

// CalculateForPeriod computes margin breakdowns for all orders within a date range,
// optionally filtered by seller. Returns breakdowns, total count, and any error.
func (s *Service) CalculateForPeriod(
	ctx context.Context,
	startDate, endDate time.Time,
	sellerID *uuid.UUID,
	page, pageSize int,
) ([]margin.MarginBreakdown, int64, error) {
	filter := order.ListFilter{
		StartDate:  &startDate,
		EndDate:    &endDate,
		SellerID:   sellerID,
		Pagination: shared.NewPagination(page, pageSize),
	}

	result, err := s.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	breakdowns := make([]margin.MarginBreakdown, 0, len(result.Items))
	for _, o := range result.Items {
		sl, err := s.sellerRepo.GetByID(ctx, o.SellerID())
		if err != nil {
			return nil, 0, err
		}

		discounts, err := s.discountRepo.ListByOrder(ctx, o.ID())
		if err != nil {
			return nil, 0, err
		}

		breakdown, err := margin.CalculateMargin(o, sl, discounts)
		if err != nil {
			return nil, 0, err
		}

		breakdowns = append(breakdowns, breakdown)
	}

	return breakdowns, result.TotalCount, nil
}
