package cloudaccount

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type SyncStatus string

const (
	SyncIdle    SyncStatus = "idle"
	SyncRunning SyncStatus = "running"
	SyncFailed  SyncStatus = "failed"
)

type CloudAccount struct {
	id           uuid.UUID
	provider     Provider
	name         string
	externalID   string
	status       Status
	syncStatus   SyncStatus
	lastSyncedAt *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func NewCloudAccount(provider Provider, name, externalID string) (*CloudAccount, error) {
	if !provider.IsValid() {
		return nil, ErrInvalidProvider
	}
	if name == "" {
		return nil, ErrAccountNameEmpty
	}
	if externalID == "" {
		return nil, ErrExternalIDEmpty
	}

	now := time.Now().UTC()
	return &CloudAccount{
		id:         uuid.New(),
		provider:   provider,
		name:       name,
		externalID: externalID,
		status:     StatusActive,
		syncStatus: SyncIdle,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func HydrateCloudAccount(
	id uuid.UUID, provider Provider, name, externalID string,
	status Status, syncStatus SyncStatus, lastSyncedAt *time.Time,
	createdAt, updatedAt time.Time,
) *CloudAccount {
	return &CloudAccount{
		id: id, provider: provider, name: name, externalID: externalID,
		status: status, syncStatus: syncStatus, lastSyncedAt: lastSyncedAt,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (a *CloudAccount) ID() uuid.UUID          { return a.id }
func (a *CloudAccount) Provider() Provider      { return a.provider }
func (a *CloudAccount) Name() string            { return a.name }
func (a *CloudAccount) ExternalID() string      { return a.externalID }
func (a *CloudAccount) Status() Status          { return a.status }
func (a *CloudAccount) SyncStatus() SyncStatus  { return a.syncStatus }
func (a *CloudAccount) LastSyncedAt() *time.Time { return a.lastSyncedAt }
func (a *CloudAccount) CreatedAt() time.Time    { return a.createdAt }
func (a *CloudAccount) UpdatedAt() time.Time    { return a.updatedAt }

func (a *CloudAccount) Rename(name string) error {
	if name == "" {
		return ErrAccountNameEmpty
	}
	a.name = name
	a.updatedAt = time.Now().UTC()
	return nil
}

func (a *CloudAccount) Deactivate() {
	a.status = StatusInactive
	a.updatedAt = time.Now().UTC()
}

func (a *CloudAccount) Activate() {
	a.status = StatusActive
	a.updatedAt = time.Now().UTC()
}

// BeginSync transitions sync state. Returns error if already syncing.
func (a *CloudAccount) BeginSync() error {
	if a.status == StatusInactive {
		return ErrAccountInactive
	}
	if a.syncStatus == SyncRunning {
		return ErrSyncAlreadyRunning
	}
	a.syncStatus = SyncRunning
	a.updatedAt = time.Now().UTC()
	return nil
}

func (a *CloudAccount) CompleteSyncSuccess() {
	now := time.Now().UTC()
	a.syncStatus = SyncIdle
	a.lastSyncedAt = &now
	a.updatedAt = now
}

func (a *CloudAccount) CompleteSyncFailure() {
	a.syncStatus = SyncFailed
	a.updatedAt = time.Now().UTC()
}
