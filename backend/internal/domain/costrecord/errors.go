package costrecord

import "errors"

var (
	ErrRecordNotFound    = errors.New("cost record not found")
	ErrServiceEmpty      = errors.New("service name cannot be empty")
	ErrInvalidAmount     = errors.New("amount must be non-negative")
	ErrInvalidUsageDate  = errors.New("usage date is required")
)
