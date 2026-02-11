package tenant

import "errors"

var (
	ErrTenantNotFound    = errors.New("tenant not found")
	ErrTenantNameEmpty   = errors.New("tenant name cannot be empty")
	ErrTenantSlugEmpty   = errors.New("tenant slug cannot be empty")
	ErrTenantSuspended   = errors.New("tenant is suspended")
	ErrTenantSlugTaken   = errors.New("tenant slug already taken")
)
