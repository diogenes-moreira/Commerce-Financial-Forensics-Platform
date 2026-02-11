package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/integration"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IntegrationRepo struct {
	db *gorm.DB
}

func NewIntegrationRepo(db *gorm.DB) *IntegrationRepo {
	return &IntegrationRepo{db: db}
}

func (r *IntegrationRepo) Create(ctx context.Context, i *integration.Integration) error {
	return r.db.WithContext(ctx).Create(mapper.IntegrationToModel(i)).Error
}

func (r *IntegrationRepo) GetByID(ctx context.Context, id uuid.UUID) (*integration.Integration, error) {
	var m model.Integration
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, integration.ErrIntegrationNotFound
		}
		return nil, err
	}
	return mapper.IntegrationToDomain(&m), nil
}

func (r *IntegrationRepo) List(ctx context.Context, filter integration.ListFilter) (*shared.PagedResult[*integration.Integration], error) {
	q := r.db.WithContext(ctx).Model(&model.Integration{})

	if filter.Platform != "" {
		q = q.Where("platform = ?", filter.Platform)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}

	var total int64
	q.Count(&total)

	var models []model.Integration
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*integration.Integration, len(models))
	for i := range models {
		items[i] = mapper.IntegrationToDomain(&models[i])
	}

	return &shared.PagedResult[*integration.Integration]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *IntegrationRepo) Update(ctx context.Context, i *integration.Integration) error {
	return r.db.WithContext(ctx).Save(mapper.IntegrationToModel(i)).Error
}

func (r *IntegrationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Integration{}).Error
}
