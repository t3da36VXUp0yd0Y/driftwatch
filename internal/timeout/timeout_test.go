package timeout

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNew_InvalidDuration(t *testing.T) {
	_, err := New(Config{Duration: 0})
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNew_Valid(t *testing.T) {
	r, err := New(Config{Duration: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRun_Success(t *testing.T) {
	r, _ := New(Config{Duration: time.Second})
	err := r.Run(context.Background(), func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PropagatesError(t *testing.T) {
	r, _ := New(Config{Duration: time.Second})
	sentinel := errors.New("boom")
	err := r.Run(context.Background(), func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestRun_Timeout(t *testing.T) {
	r, _ := New(Config{Duration: 20 * time.Millisecond})
	err := r.Run(context.Background(), func(ctx context.Context) error {
		select {
		case <-time.After(time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Duration != 30*time.Second {
		t.Fatalf("expected 30s, got %v", cfg.Duration)
	}
}
