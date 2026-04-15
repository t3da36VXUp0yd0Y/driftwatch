// Package retry provides a simple retry mechanism with configurable
// attempts and backoff for use when polling Docker or external services.
package retry

import (
	"context"
	"fmt"
	"time"
)

// Config holds the retry configuration.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// Delay is the wait duration between attempts.
	Delay time.Duration
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
	}
}

// Do executes fn up to cfg.MaxAttempts times, returning nil on the first
// success. If all attempts fail the last error is returned wrapped with
// attempt metadata. The context is checked before every attempt.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry cancelled before attempt %d: %w", attempt, err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt < cfg.MaxAttempts {
			select {
			case <-time.After(cfg.Delay):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled during backoff after attempt %d: %w", attempt, ctx.Err())
			}
		}
	}

	return fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}
