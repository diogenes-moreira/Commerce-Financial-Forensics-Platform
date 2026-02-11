package budget

import "errors"

var (
	ErrBudgetNotFound   = errors.New("budget not found")
	ErrBudgetNameEmpty  = errors.New("budget name cannot be empty")
	ErrInvalidPeriod    = errors.New("period start must be before period end")
	ErrBudgetExceeded   = errors.New("budget has been exceeded")
)
