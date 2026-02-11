package order

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrExternalIDEmpty    = errors.New("order external ID cannot be empty")
	ErrInvalidTransition  = errors.New("invalid order status transition")
	ErrNoItems            = errors.New("order must have at least one item")
	ErrInvalidQuantity    = errors.New("quantity must be greater than zero")
	ErrItemNotFound       = errors.New("order item not found")
)
