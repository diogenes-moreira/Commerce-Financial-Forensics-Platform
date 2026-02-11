package customer

import "errors"

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNameEmpty        = errors.New("customer name cannot be empty")
)
