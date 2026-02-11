package integration

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationStatus string

const (
	StatusActive   IntegrationStatus = "active"
	StatusInactive IntegrationStatus = "inactive"
	StatusError    IntegrationStatus = "error"
)

var allowedPlatforms = map[string]bool{
	"woocommerce": true,
	"shopify":     true,
	"medusajs":    true,
	"custom":      true,
}

type Integration struct {
	id           uuid.UUID
	platform     string
	name         string
	config       map[string]string
	status       IntegrationStatus
	syncStatus   string
	lastSyncedAt *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewIntegration(platform, name string, config map[string]string) (*Integration, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}
	if !allowedPlatforms[platform] {
		return nil, ErrInvalidPlatform
	}
	if config == nil {
		config = make(map[string]string)
	}

	now := time.Now().UTC()
	return &Integration{
		id:         uuid.New(),
		platform:   platform,
		name:       name,
		config:     config,
		status:     StatusActive,
		syncStatus: "idle",
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func HydrateIntegration(
	id uuid.UUID, platform, name string, config map[string]string,
	status IntegrationStatus, syncStatus string, lastSyncedAt *time.Time,
	createdAt, updatedAt time.Time,
) *Integration {
	if config == nil {
		config = make(map[string]string)
	}
	return &Integration{
		id: id, platform: platform, name: name, config: config,
		status: status, syncStatus: syncStatus, lastSyncedAt: lastSyncedAt,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (i *Integration) ID() uuid.UUID             { return i.id }
func (i *Integration) Platform() string           { return i.platform }
func (i *Integration) Name() string               { return i.name }
func (i *Integration) Config() map[string]string  { return i.config }
func (i *Integration) Status() IntegrationStatus  { return i.status }
func (i *Integration) SyncStatus() string         { return i.syncStatus }
func (i *Integration) LastSyncedAt() *time.Time   { return i.lastSyncedAt }
func (i *Integration) CreatedAt() time.Time       { return i.createdAt }
func (i *Integration) UpdatedAt() time.Time       { return i.updatedAt }

func (i *Integration) Activate() {
	i.status = StatusActive
	i.updatedAt = time.Now().UTC()
}

func (i *Integration) Deactivate() {
	i.status = StatusInactive
	i.updatedAt = time.Now().UTC()
}

func (i *Integration) UpdateDetails(name string, config map[string]string) error {
	if name != "" {
		i.name = name
	}
	if config != nil {
		i.config = config
	}
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Integration) BeginSync() error {
	if i.status == StatusInactive {
		return ErrIntegrationInactive
	}
	if i.syncStatus == "running" {
		return ErrAlreadySyncing
	}
	i.syncStatus = "running"
	i.updatedAt = time.Now().UTC()
	return nil
}

func (i *Integration) CompleteSync() {
	now := time.Now().UTC()
	i.syncStatus = "idle"
	i.lastSyncedAt = &now
	i.updatedAt = now
}

func (i *Integration) FailSync(reason string) {
	i.syncStatus = "failed"
	i.status = StatusError
	i.updatedAt = time.Now().UTC()
}
