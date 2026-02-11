package anomaly

import "errors"

var (
	ErrAnomalyNotFound     = errors.New("anomaly not found")
	ErrAlreadyResolved     = errors.New("anomaly is already resolved")
)
