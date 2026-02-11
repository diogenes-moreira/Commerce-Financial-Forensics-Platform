package ledger

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/ledger"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo ledger.Repository
}

func NewService(repo ledger.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	orderID *uuid.UUID, accountCode string, side ledger.Side,
	amount shared.Money, description, refType string,
	refID uuid.UUID, effectiveDate time.Time,
) (*ledger.LedgerEntry, error) {
	entry, err := ledger.NewLedgerEntry(
		orderID, accountCode, side, amount,
		description, refType, refID, effectiveDate,
	)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*ledger.LedgerEntry, error) {
	return s.repo.ListByOrder(ctx, orderID)
}

func (s *Service) ListByAccount(ctx context.Context, code string, dateRange ledger.DateRange, pagination shared.Pagination) (*shared.PagedResult[*ledger.LedgerEntry], error) {
	return s.repo.ListByAccount(ctx, code, dateRange, pagination)
}

func (s *Service) BalanceByAccount(ctx context.Context, code string, dateRange ledger.DateRange) (int64, error) {
	return s.repo.BalanceByAccount(ctx, code, dateRange)
}
