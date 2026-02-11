package forensicevent

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/forensicevent"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	repo forensicevent.Repository
}

func NewService(repo forensicevent.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	entityType string, entityID uuid.UUID, eventType string,
	effectiveTime time.Time, payload map[string]any, actorID string,
) (*forensicevent.ForensicEvent, error) {
	event, err := forensicevent.NewForensicEvent(
		entityType, entityID, eventType,
		effectiveTime, payload, actorID,
	)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*forensicevent.ForensicEvent, error) {
	return s.repo.ListByEntity(ctx, entityType, entityID)
}

func (s *Service) ListByType(ctx context.Context, eventType string, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*forensicevent.ForensicEvent], error) {
	return s.repo.ListByType(ctx, eventType, start, end, pagination)
}

func (s *Service) ListByDateRange(ctx context.Context, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*forensicevent.ForensicEvent], error) {
	return s.repo.ListByDateRange(ctx, start, end, pagination)
}
