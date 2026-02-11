package forensicevent

import "errors"

var (
	ErrEventNotFound     = errors.New("forensic event not found")
	ErrEntityTypeEmpty   = errors.New("entity type cannot be empty")
	ErrEventTypeEmpty    = errors.New("event type cannot be empty")
)
