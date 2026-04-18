package metric

import (
	"testing"
)

func TestNew_StartsAtZero(t *testing.T) {
	c := New()
	if got := c.Get("foo"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestInc_IncrementsCounter(t *testing.T) {
	c := New()
	c.Inc("hits")
	c.Inc("hits")
	if got := c.Get("hits"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestAdd_AddsValue(t *testing.T) {
	c := New()
	c.Add("bytes", 512)
	c.Add("bytes", 256)
	if got := c.Get("bytes"); got != 768 {
		t.Fatalf("expected 768, got %d", got)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	c := New()
	c.Inc("a")
	c.Inc("b")
	snap, _ := c.Snapshot()
	if snap["a"] != 1 || snap["b"] != 1 {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
	// mutating snapshot must not affect counter
	snap["a"] = 99
	if c.Get("a") != 1 {
		t.Fatal("snapshot mutation affected counter")
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	c := New()
	c.Inc("x")
	c.Reset()
	if got := c.Get("x"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
