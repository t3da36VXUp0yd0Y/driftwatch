package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Remaining() != 3 {
		t.Fatalf("expected 3 tokens, got %d", l.Remaining())
	}
}

func TestAllow_ConsumesTokens(t *testing.T) {
	l, _ := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("expected allow on call %d, got %v", i+1, err)
		}
	}
	if l.Remaining() != 0 {
		t.Fatalf("expected 0 tokens remaining")
	}
}

func TestAllow_RateLimitedWhenEmpty(t *testing.T) {
	l, _ := New(2, time.Minute)
	l.Allow()
	l.Allow()
	if err := l.Allow(); err != ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestAllow_RefillsAfterWindow(t *testing.T) {
	l, _ := New(2, 50*time.Millisecond)
	l.Allow()
	l.Allow()
	if err := l.Allow(); err != ErrRateLimited {
		t.Fatal("should be rate limited before window expires")
	}
	time.Sleep(60 * time.Millisecond)
	if err := l.Allow(); err != nil {
		t.Fatalf("expected allow after window refill, got %v", err)
	}
	if l.Remaining() != 1 {
		t.Fatalf("expected 1 token remaining after one call post-refill")
	}
}
