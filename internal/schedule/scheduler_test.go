package schedule_test

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/schedule"
)

type mockRunner struct {
	calls atomic.Int32
	err   error
}

func (m *mockRunner) Run(_ context.Context) error {
	m.calls.Add(1)
	return m.err
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := schedule.New(&mockRunner{}, 0, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
}

func TestNew_NilRunner(t *testing.T) {
	_, err := schedule.New(nil, time.Second, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for nil runner, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := schedule.New(&mockRunner{}, time.Second, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

func TestStart_RunsImmediately(t *testing.T) {
	r := &mockRunner{}
	var buf bytes.Buffer
	s, _ := schedule.New(r, 10*time.Second, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately after first tick

	// Start will run once then return ctx.Err
	err := s.Start(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if r.calls.Load() < 1 {
		t.Fatalf("expected at least 1 call, got %d", r.calls.Load())
	}
}

func TestStart_RunnerError(t *testing.T) {
	runErr := errors.New("runner failed")
	r := &mockRunner{err: runErr}
	var buf bytes.Buffer
	s, _ := schedule.New(r, time.Second, &buf)

	ctx := context.Background()
	err := s.Start(ctx)
	if !errors.Is(err, runErr) {
		t.Fatalf("expected runner error, got %v", err)
	}
}

func TestStart_TicksMultipleTimes(t *testing.T) {
	r := &mockRunner{}
	var buf bytes.Buffer
	s, _ := schedule.New(r, 20*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()

	err := s.Start(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	// Should have run at t=0, ~20ms, ~40ms, ~60ms => at least 3 calls
	if r.calls.Load() < 3 {
		t.Fatalf("expected at least 3 calls, got %d", r.calls.Load())
	}
}
