package watch

import "errors"

// ErrInvalidInterval is returned when the poll interval is zero or negative.
var ErrInvalidInterval = errors.New("watch: interval must be greater than zero")

// ErrNilCallback is returned when no onChange callback is provided.
var ErrNilCallback = errors.New("watch: onChange callback must not be nil")
