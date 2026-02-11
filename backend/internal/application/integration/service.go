package integration

import (
	"context"
	"fmt"

	"github.com/diogenes/costforensics/backend/internal/domain/integration"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/diogenes/costforensics/backend/internal/platform/connector"
	"github.com/google/uuid"
)

type Service struct {
	repo integration.Repository
}

func NewService(repo integration.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, platform, name string, config map[string]string) (*integration.Integration, error) {
	i, err := integration.NewIntegration(platform, name, config)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*integration.Integration, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter integration.ListFilter) (*shared.PagedResult[*integration.Integration], error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name string, config map[string]string) (*integration.Integration, error) {
	i, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.UpdateDetails(name, config); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) TriggerSync(ctx context.Context, id uuid.UUID) (*integration.Integration, error) {
	i, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.BeginSync(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, i); err != nil {
		return nil, err
	}

	// For now, immediately mark sync as completed.
	// Actual connector logic will be wired in later.
	i.CompleteSync()
	if err := s.repo.Update(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) GetConnector(i *integration.Integration) (connector.Connector, error) {
	switch i.Platform() {
	case "woocommerce":
		return connector.NewWooCommerceConnector(i.Config()), nil
	case "shopify":
		return connector.NewShopifyConnector(i.Config()), nil
	case "medusajs":
		return connector.NewMedusaJSConnector(i.Config()), nil
	case "custom":
		return nil, fmt.Errorf("custom connector requires external configuration")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", i.Platform())
	}
}
