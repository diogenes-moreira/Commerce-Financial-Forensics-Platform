package ledger

import "errors"

var (
	ErrEntryNotFound      = errors.New("ledger entry not found")
	ErrInvalidAccountCode = errors.New("account code cannot be empty")
	ErrInvalidSide        = errors.New("side must be 'debit' or 'credit'")
)
