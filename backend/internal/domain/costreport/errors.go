package costreport

import "errors"

var (
	ErrReportNotFound  = errors.New("cost report not found")
	ErrReportNameEmpty = errors.New("report name cannot be empty")
	ErrInvalidPeriod   = errors.New("period start must be before period end")
)
