package mapper

import (
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
)

func TenantToModel(t *tenant.Tenant) *model.Tenant {
	return &model.Tenant{
		ID:        t.ID(),
		Name:      t.Name(),
		Slug:      t.Slug(),
		DBName:    t.DBName(),
		Status:    string(t.Status()),
		CreatedAt: t.CreatedAt(),
		UpdatedAt: t.UpdatedAt(),
	}
}

func TenantToDomain(m *model.Tenant) *tenant.Tenant {
	return tenant.Hydrate(
		m.ID,
		m.Name,
		m.Slug,
		m.DBName,
		tenant.Status(m.Status),
		m.CreatedAt,
		m.UpdatedAt,
	)
}
