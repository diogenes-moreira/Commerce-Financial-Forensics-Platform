package product

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrSKUEmpty         = errors.New("product SKU cannot be empty")
	ErrNameEmpty        = errors.New("product name cannot be empty")
	ErrPriceBelowCost   = errors.New("unit price must be greater than or equal to unit cost")
)
