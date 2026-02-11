package discount

import "errors"

var (
	ErrRuleNotFound       = errors.New("promotion rule not found")
	ErrApplicationNotFound = errors.New("discount application not found")
	ErrNameEmpty          = errors.New("promotion name cannot be empty")
	ErrRuleExpired        = errors.New("promotion rule has expired")
	ErrMaxUsageReached    = errors.New("promotion maximum usage reached")
)
