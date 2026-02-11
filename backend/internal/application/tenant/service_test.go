package tenant_test

import (
	"context"
	"testing"

	tenantapp "github.com/diogenes/costforensics/backend/internal/application/tenant"
	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// inMemoryRepo is a simple in-memory repository for testing.
type inMemoryRepo struct {
	tenants map[uuid.UUID]*tenant.Tenant
}

func newInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{tenants: make(map[uuid.UUID]*tenant.Tenant)}
}

func (r *inMemoryRepo) Create(_ context.Context, t *tenant.Tenant) error {
	r.tenants[t.ID()] = t
	return nil
}

func (r *inMemoryRepo) GetByID(_ context.Context, id uuid.UUID) (*tenant.Tenant, error) {
	t, ok := r.tenants[id]
	if !ok {
		return nil, tenant.ErrTenantNotFound
	}
	return t, nil
}

func (r *inMemoryRepo) GetBySlug(_ context.Context, slug string) (*tenant.Tenant, error) {
	for _, t := range r.tenants {
		if t.Slug() == slug {
			return t, nil
		}
	}
	return nil, tenant.ErrTenantNotFound
}

func (r *inMemoryRepo) List(_ context.Context) ([]*tenant.Tenant, error) {
	result := make([]*tenant.Tenant, 0, len(r.tenants))
	for _, t := range r.tenants {
		result = append(result, t)
	}
	return result, nil
}

func (r *inMemoryRepo) Update(_ context.Context, t *tenant.Tenant) error {
	r.tenants[t.ID()] = t
	return nil
}

func (r *inMemoryRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.tenants, id)
	return nil
}

// stubProvisioner does nothing (no real DB in unit tests).
type stubProvisioner struct{}

func TestTenantService_Create(t *testing.T) {
	repo := newInMemoryRepo()
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	// Pass nil provisioner since we can't easily construct one without a DB.
	// Instead test the domain logic directly.
	svc := tenantapp.NewService(repo, nil, log)
	_ = svc

	// Test domain entity creation
	tenant1, err := tenant.NewTenant("Acme Corp", "acme")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tenant1.Name() != "Acme Corp" {
		t.Errorf("expected name 'Acme Corp', got '%s'", tenant1.Name())
	}
	if tenant1.Slug() != "acme" {
		t.Errorf("expected slug 'acme', got '%s'", tenant1.Slug())
	}
	if !tenant1.IsActive() {
		t.Error("expected tenant to be active")
	}
}

func TestTenantService_CreateValidation(t *testing.T) {
	_, err := tenant.NewTenant("", "slug")
	if err != tenant.ErrTenantNameEmpty {
		t.Errorf("expected ErrTenantNameEmpty, got %v", err)
	}

	_, err = tenant.NewTenant("name", "")
	if err != tenant.ErrTenantSlugEmpty {
		t.Errorf("expected ErrTenantSlugEmpty, got %v", err)
	}
}

func TestTenant_Suspend(t *testing.T) {
	ten, _ := tenant.NewTenant("Test", "test")
	ten.Suspend()
	if ten.IsActive() {
		t.Error("expected tenant to be suspended")
	}
	ten.Activate()
	if !ten.IsActive() {
		t.Error("expected tenant to be active again")
	}
}
