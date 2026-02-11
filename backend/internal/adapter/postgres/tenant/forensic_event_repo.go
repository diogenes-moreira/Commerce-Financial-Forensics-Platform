package tenant

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/forensicevent"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForensicEventRepo struct {
	db *gorm.DB
}

func NewForensicEventRepo(db *gorm.DB) *ForensicEventRepo {
	return &ForensicEventRepo{db: db}
}

func (r *ForensicEventRepo) Create(ctx context.Context, e *forensicevent.ForensicEvent) error {
	return r.db.WithContext(ctx).Create(mapper.ForensicEventToModel(e)).Error
}

func (r *ForensicEventRepo) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*forensicevent.ForensicEvent, error) {
	var models []model.ForensicEvent
	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("event_time DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]*forensicevent.ForensicEvent, len(models))
	for i := range models {
		items[i] = mapper.ForensicEventToDomain(&models[i])
	}
	return items, nil
}

func (r *ForensicEventRepo) ListByType(ctx context.Context, eventType string, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*forensicevent.ForensicEvent], error) {
	q := r.db.WithContext(ctx).Model(&model.ForensicEvent{}).Where("event_type = ?", eventType)

	if !start.IsZero() {
		q = q.Where("event_time >= ?", start)
	}
	if !end.IsZero() {
		q = q.Where("event_time <= ?", end)
	}

	var total int64
	q.Count(&total)

	var models []model.ForensicEvent
	if err := q.Offset(pagination.Offset()).Limit(pagination.PageSize()).
		Order("event_time DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*forensicevent.ForensicEvent, len(models))
	for i := range models {
		items[i] = mapper.ForensicEventToDomain(&models[i])
	}

	return &shared.PagedResult[*forensicevent.ForensicEvent]{
		Items:      items,
		TotalCount: total,
		Page:       pagination.Page(),
		PageSize:   pagination.PageSize(),
	}, nil
}

func (r *ForensicEventRepo) ListByDateRange(ctx context.Context, start, end time.Time, pagination shared.Pagination) (*shared.PagedResult[*forensicevent.ForensicEvent], error) {
	q := r.db.WithContext(ctx).Model(&model.ForensicEvent{})

	if !start.IsZero() {
		q = q.Where("event_time >= ?", start)
	}
	if !end.IsZero() {
		q = q.Where("event_time <= ?", end)
	}

	var total int64
	q.Count(&total)

	var models []model.ForensicEvent
	if err := q.Offset(pagination.Offset()).Limit(pagination.PageSize()).
		Order("event_time DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*forensicevent.ForensicEvent, len(models))
	for i := range models {
		items[i] = mapper.ForensicEventToDomain(&models[i])
	}

	return &shared.PagedResult[*forensicevent.ForensicEvent]{
		Items:      items,
		TotalCount: total,
		Page:       pagination.Page(),
		PageSize:   pagination.PageSize(),
	}, nil
}
