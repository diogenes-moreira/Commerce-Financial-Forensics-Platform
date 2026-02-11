package importjob

import "errors"

var (
	ErrNameEmpty         = errors.New("import job name is required")
	ErrInvalidEntityType = errors.New("invalid entity type")
	ErrJobNotFound       = errors.New("import job not found")
	ErrAlreadyProcessing = errors.New("import job is already processing")
)
