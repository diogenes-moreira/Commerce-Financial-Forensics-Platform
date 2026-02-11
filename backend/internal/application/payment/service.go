package payment

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/payment"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo payment.Repository
}

func NewService(repo payment.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	orderID uuid.UUID, direction payment.Direction, counterpartyID uuid.UUID,
	amount shared.Money, method, externalRef string,
) (*payment.Payment, error) {
	p, err := payment.NewPayment(
		orderID, direction, counterpartyID,
		amount, method, externalRef,
	)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*payment.Payment, error) {
	return s.repo.ListByOrder(ctx, orderID)
}

func (s *Service) List(ctx context.Context, filter payment.ListFilter) (*shared.PagedResult[*payment.Payment], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) MarkProcessed(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.MarkProcessed(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) MarkFailed(ctx context.Context, id uuid.UUID, reason string) (*payment.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.MarkFailed(reason); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) MarkRefunded(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := p.MarkRefunded(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
