// Package schedule provides periodic drift-check scheduling for driftwatch.
package schedule

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Runner is the interface satisfied by runner.Runner.
type Runner interface {
	Run(ctx context.Context) error
}

// Scheduler runs a Runner on a fixed interval until the context is cancelled.
type Scheduler struct {
	runner   Runner
	interval time.Duration
	out      io.Writer
}

// New returns a Scheduler that will invoke r every interval.
// interval must be greater than zero.
func New(r Runner, interval time.Duration, out io.Writer) (*Scheduler, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("schedule: interval must be greater than zero, got %s", interval)
	}
	if r == nil {
		return nil, fmt.Errorf("schedule: runner must not be nil")
	}
	return &Scheduler{runner: r, interval: interval, out: out}, nil
}

// Start blocks and runs the drift check immediately, then on every interval.
// It returns when ctx is cancelled, propagating ctx.Err().
func (s *Scheduler) Start(ctx context.Context) error {
	if err := s.tick(ctx); err != nil {
		return err
	}
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.C:
			fmt.Fprintf(s.out, "[schedule] running drift check at %s\n", t.Format(time.RFC3339))
			if err := s.tick(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) error {
	return s.runner.Run(ctx)
}
