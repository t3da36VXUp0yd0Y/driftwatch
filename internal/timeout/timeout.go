package timeout

import (
	"context"
	"fmt"
	"time"
)

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = fmt.Errorf("operation timed out")

// Config holds timeout settings.
type Config struct {
	Duration time.Duration
}

// DefaultConfig returns a Config with a sensible default.
func DefaultConfig() Config {
	return Config{Duration: 30 * time.Second}
}

// Runner wraps a function with a timeout context.
type Runner struct {
	cfg Config
}

// New creates a Runner. Returns an error if Duration is zero or negative.
func New(cfg Config) (*Runner, error) {
	if cfg.Duration <= 0 {
		return nil, fmt.Errorf("timeout: duration must be positive, got %v", cfg.Duration)
	}
	return &Runner{cfg: cfg}, nil
}

// Run executes fn within the configured timeout.
// Returns ErrTimeout if the deadline is exceeded.
func (r *Runner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Duration)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		ch <- fn(ctx)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return ErrTimeout
		}
		return ctx.Err()
	}
}
