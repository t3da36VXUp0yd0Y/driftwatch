package throttle_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/throttle"
)

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := throttle.New(0)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
	_, err = throttle.New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative cooldown")
	}
}

func TestNew_Valid(t *testing.T) {
	th, err := throttle.New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}

func TestAllow_FirstCallAllowed(t *testing.T) {
	th, _ := throttle.New(time.Second)
	if !th.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SuppressedWithinCooldown(t *testing.T) {
	th, _ := throttle.New(time.Hour)
	th.Allow("svc-a")
	if th.Allow("svc-a") {
		t.Fatal("expected second call to be suppressed within cooldown")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	th, _ := throttle.New(10 * time.Millisecond)
	th.Allow("svc-a")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("svc-a") {
		t.Fatal("expected call to be allowed after cooldown expires")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th, _ := throttle.New(time.Hour)
	th.Allow("svc-a")
	if !th.Allow("svc-b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	th, _ := throttle.New(time.Hour)
	th.Allow("svc-a")
	th.Reset("svc-a")
	if !th.Allow("svc-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAll(t *testing.T) {
	th, _ := throttle.New(time.Hour)
	th.Allow("svc-a")
	th.Allow("svc-b")
	th.ResetAll()
	if !th.Allow("svc-a") || !th.Allow("svc-b") {
		t.Fatal("expected all keys allowed after ResetAll")
	}
}
