package tenant

import (
	"context"

	"github.com/diogenes/costforensics/backend/internal/domain/tenant"
	"github.com/diogenes/costforensics/backend/internal/platform/database"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo        tenant.Repository
	provisioner *database.Provisioner
	log         *logrus.Logger
}

func NewService(repo tenant.Repository, provisioner *database.Provisioner, log *logrus.Logger) *Service {
	return &Service{repo: repo, provisioner: provisioner, log: log}
}

func (s *Service) Create(ctx context.Context, name, slug string) (*tenant.Tenant, error) {
	// Check slug uniqueness
	if _, err := s.repo.GetBySlug(ctx, slug); err == nil {
		return nil, tenant.ErrTenantSlugTaken
	}

	t, err := tenant.NewTenant(name, slug)
	if err != nil {
		return nil, err
	}

	// Provision tenant database
	if err := s.provisioner.ProvisionTenantDB(t.DBName()); err != nil {
		return nil, err
	}

	// Register in system DB
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	s.log.WithField("tenant_id", t.ID()).Info("Tenant created")
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*tenant.Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*tenant.Tenant, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name string) (*tenant.Tenant, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := t.Rename(name); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
