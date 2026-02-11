package shared

import (
	"errors"
	"time"
)

var ErrInvalidDateRange = errors.New("start date must be before end date")

// DateRange is an immutable value object representing a time period.
type DateRange struct {
	start time.Time
	end   time.Time
}

func NewDateRange(start, end time.Time) (DateRange, error) {
	if !start.Before(end) {
		return DateRange{}, ErrInvalidDateRange
	}
	return DateRange{start: start, end: end}, nil
}

func (d DateRange) Start() time.Time { return d.start }
func (d DateRange) End() time.Time   { return d.end }

func (d DateRange) Contains(t time.Time) bool {
	return !t.Before(d.start) && !t.After(d.end)
}

func (d DateRange) Overlaps(other DateRange) bool {
	return d.start.Before(other.end) && other.start.Before(d.end)
}

func (d DateRange) Days() int {
	return int(d.end.Sub(d.start).Hours() / 24)
}
