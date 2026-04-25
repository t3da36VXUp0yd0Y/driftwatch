package window

import "errors"

// ErrInvalidDuration is returned when a non-positive duration is provided to New.
var ErrInvalidDuration = errors.New("window: duration must be greater than zero")
