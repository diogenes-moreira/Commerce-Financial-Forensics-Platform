package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImportJobRepo struct {
	db *gorm.DB
}

func NewImportJobRepo(db *gorm.DB) *ImportJobRepo {
	return &ImportJobRepo{db: db}
}

func (r *ImportJobRepo) Create(ctx context.Context, j *importjob.ImportJob) error {
	return r.db.WithContext(ctx).Create(mapper.ImportJobToModel(j)).Error
}

func (r *ImportJobRepo) GetByID(ctx context.Context, id uuid.UUID) (*importjob.ImportJob, error) {
	var m model.ImportJob
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, importjob.ErrJobNotFound
		}
		return nil, err
	}
	return mapper.ImportJobToDomain(&m), nil
}

func (r *ImportJobRepo) List(ctx context.Context, filter importjob.ListFilter) (*shared.PagedResult[*importjob.ImportJob], error) {
	q := r.db.WithContext(ctx).Model(&model.ImportJob{})

	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.EntityType != "" {
		q = q.Where("entity_type = ?", filter.EntityType)
	}

	var total int64
	q.Count(&total)

	var models []model.ImportJob
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*importjob.ImportJob, len(models))
	for i := range models {
		items[i] = mapper.ImportJobToDomain(&models[i])
	}

	return &shared.PagedResult[*importjob.ImportJob]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *ImportJobRepo) Update(ctx context.Context, j *importjob.ImportJob) error {
	return r.db.WithContext(ctx).Save(mapper.ImportJobToModel(j)).Error
}
