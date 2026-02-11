package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
)

// Tenant is a rich domain entity. Fields are unexported; mutation via methods.
type Tenant struct {
	id        uuid.UUID
	name      string
	slug      string
	dbName    string
	status    Status
	createdAt time.Time
	updatedAt time.Time
}

func NewTenant(name, slug string) (*Tenant, error) {
	if name == "" {
		return nil, ErrTenantNameEmpty
	}
	if slug == "" {
		return nil, ErrTenantSlugEmpty
	}

	now := time.Now().UTC()
	id := uuid.New()
	return &Tenant{
		id:        id,
		name:      name,
		slug:      slug,
		dbName:    "cf_tenant_" + id.String()[:8],
		status:    StatusActive,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Hydrate reconstructs a Tenant from persistence. Not for business logic creation.
func Hydrate(id uuid.UUID, name, slug, dbName string, status Status, createdAt, updatedAt time.Time) *Tenant {
	return &Tenant{
		id:        id,
		name:      name,
		slug:      slug,
		dbName:    dbName,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *Tenant) ID() uuid.UUID       { return t.id }
func (t *Tenant) Name() string         { return t.name }
func (t *Tenant) Slug() string         { return t.slug }
func (t *Tenant) DBName() string       { return t.dbName }
func (t *Tenant) Status() Status       { return t.status }
func (t *Tenant) CreatedAt() time.Time { return t.createdAt }
func (t *Tenant) UpdatedAt() time.Time { return t.updatedAt }
func (t *Tenant) IsActive() bool       { return t.status == StatusActive }

func (t *Tenant) Rename(name string) error {
	if name == "" {
		return ErrTenantNameEmpty
	}
	t.name = name
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Tenant) Suspend() error {
	t.status = StatusSuspended
	t.updatedAt = time.Now().UTC()
	return nil
}

func (t *Tenant) Activate() {
	t.status = StatusActive
	t.updatedAt = time.Now().UTC()
}
