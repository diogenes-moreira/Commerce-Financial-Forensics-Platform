package system

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(ctx context.Context, t *tenant.Tenant) error {
	m := mapper.TenantToModel(t)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *TenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*tenant.Tenant, error) {
	var m model.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, tenant.ErrTenantNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapper.TenantToDomain(&m), nil
}

func (r *TenantRepo) GetBySlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	var m model.Tenant
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, tenant.ErrTenantNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapper.TenantToDomain(&m), nil
}

func (r *TenantRepo) List(ctx context.Context) ([]*tenant.Tenant, error) {
	var models []model.Tenant
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	tenants := make([]*tenant.Tenant, len(models))
	for i := range models {
		tenants[i] = mapper.TenantToDomain(&models[i])
	}
	return tenants, nil
}

func (r *TenantRepo) Update(ctx context.Context, t *tenant.Tenant) error {
	m := mapper.TenantToModel(t)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TenantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tenant{}).Error
}
