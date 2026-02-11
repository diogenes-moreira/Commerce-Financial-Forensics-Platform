package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/integration"
)

func IntegrationToModel(i *integration.Integration) *model.Integration {
	configJSON, _ := json.Marshal(i.Config())
	return &model.Integration{
		ID:           i.ID(),
		Platform:     i.Platform(),
		Name:         i.Name(),
		Config:       string(configJSON),
		Status:       string(i.Status()),
		SyncStatus:   i.SyncStatus(),
		LastSyncedAt: i.LastSyncedAt(),
		CreatedAt:    i.CreatedAt(),
		UpdatedAt:    i.UpdatedAt(),
	}
}

func IntegrationToDomain(m *model.Integration) *integration.Integration {
	var config map[string]string
	_ = json.Unmarshal([]byte(m.Config), &config)
	if config == nil {
		config = make(map[string]string)
	}
	return integration.HydrateIntegration(
		m.ID, m.Platform, m.Name, config,
		integration.IntegrationStatus(m.Status), m.SyncStatus, m.LastSyncedAt,
		m.CreatedAt, m.UpdatedAt,
	)
}
