package integration

import "errors"

var (
	ErrNameEmpty            = errors.New("integration name cannot be empty")
	ErrInvalidPlatform      = errors.New("invalid platform; must be woocommerce, shopify, medusajs, or custom")
	ErrIntegrationNotFound  = errors.New("integration not found")
	ErrAlreadySyncing       = errors.New("sync is already running")
	ErrIntegrationInactive  = errors.New("integration is inactive")
)
