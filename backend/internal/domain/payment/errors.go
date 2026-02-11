package payment

import "errors"

var (
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrInvalidDirection   = errors.New("direction must be 'inbound' or 'outbound'")
	ErrInvalidStatus      = errors.New("invalid payment status transition")
	ErrAmountRequired     = errors.New("payment amount is required")
)
