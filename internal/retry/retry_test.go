package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/retry"
)

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 3, Delay: 0}, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	sentinel := errors.New("temporary error")
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 3, Delay: 0}, func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AllAttemptsFail(t *testing.T) {
	sentinel := errors.New("persistent error")
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 2, Delay: 0}, func() error {
		return sentinel
	})
	if err == nil {
		t.Fatal("expected error when all attempts fail")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped sentinel, got %v", err)
	}
}

func TestDo_ContextCancelledBeforeAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, retry.Config{MaxAttempts: 3, Delay: 0}, func() error {
		calls++
		return nil
	})
	if err == nil {
		t.Fatal("expected error when context is cancelled")
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls with cancelled context, got %d", calls)
	}
}

func TestDo_ContextCancelledDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	err := retry.Do(ctx, retry.Config{MaxAttempts: 3, Delay: 5 * time.Second}, func() error {
		calls++
		cancel() // cancel after first failure to trigger backoff cancellation
		return errors.New("fail")
	})
	if err == nil {
		t.Fatal("expected error when context cancelled during backoff")
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call before backoff cancellation, got %d", calls)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay=500ms, got %v", cfg.Delay)
	}
}
