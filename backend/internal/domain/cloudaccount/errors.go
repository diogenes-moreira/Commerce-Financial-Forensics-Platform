package cloudaccount

import "errors"

var (
	ErrAccountNotFound     = errors.New("cloud account not found")
	ErrAccountNameEmpty    = errors.New("account name cannot be empty")
	ErrInvalidProvider     = errors.New("invalid cloud provider")
	ErrExternalIDEmpty     = errors.New("external id cannot be empty")
	ErrSyncAlreadyRunning  = errors.New("sync is already running")
	ErrAccountInactive     = errors.New("cloud account is inactive")
)
