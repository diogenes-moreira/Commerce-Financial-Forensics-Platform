package response

import (
	"github.com/diogenes/costforensics/backend/internal/domain/integration"
)

type Integration struct {
	ID           string            `json:"id"`
	Platform     string            `json:"platform"`
	Name         string            `json:"name"`
	Config       map[string]string `json:"config"`
	Status       string            `json:"status"`
	SyncStatus   string            `json:"sync_status"`
	LastSyncedAt *string           `json:"last_synced_at"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

func IntegrationFromDomain(i *integration.Integration) Integration {
	var lastSynced *string
	if i.LastSyncedAt() != nil {
		s := i.LastSyncedAt().Format("2006-01-02T15:04:05Z")
		lastSynced = &s
	}
	return Integration{
		ID:           i.ID().String(),
		Platform:     i.Platform(),
		Name:         i.Name(),
		Config:       i.Config(),
		Status:       string(i.Status()),
		SyncStatus:   i.SyncStatus(),
		LastSyncedAt: lastSynced,
		CreatedAt:    i.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    i.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}
}
