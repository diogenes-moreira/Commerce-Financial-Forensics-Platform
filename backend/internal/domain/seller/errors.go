package seller

import "errors"

var (
	ErrSellerNotFound      = errors.New("seller not found")
	ErrNameEmpty           = errors.New("seller name cannot be empty")
	ErrInvalidCommission   = errors.New("commission must be between 0 and 100")
)
