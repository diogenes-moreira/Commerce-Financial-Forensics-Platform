package exchangerate

import "errors"

var (
	ErrRateNotFound     = errors.New("exchange rate not found")
	ErrInvalidPair      = errors.New("currency pair cannot be empty")
	ErrInvalidRate      = errors.New("exchange rate must be positive")
)
