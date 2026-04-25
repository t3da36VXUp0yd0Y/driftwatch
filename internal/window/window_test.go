package window_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/window"
)

func makeResults(drifted bool) []drift.Result {
	return []drift.Result{
		{Service: "svc-a", Drifted: drifted},
		{Service: "svc-b", Drifted: false},
	}
}

func TestNew_InvalidDuration(t *testing.T) {
	_, err := window.New(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
	_, err = window.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNew_Valid(t *testing.T) {
	w, err := window.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestAdd_And_Entries(t *testing.T) {
	w, _ := window.New(time.Minute)
	w.Add(makeResults(false))
	w.Add(makeResults(true))

	entries := w.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestDriftedCount_CountsOnlyDrifted(t *testing.T) {
	w, _ := window.New(time.Minute)
	w.Add(makeResults(true))  // 1 drifted
	w.Add(makeResults(false)) // 0 drifted
	w.Add(makeResults(true))  // 1 drifted

	got := w.DriftedCount()
	if got != 2 {
		t.Fatalf("expected 2 drifted, got %d", got)
	}
}

func TestEntries_EvictsExpired(t *testing.T) {
	w, _ := window.New(50 * time.Millisecond)
	w.Add(makeResults(true))
	time.Sleep(80 * time.Millisecond)
	w.Add(makeResults(false))

	entries := w.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", len(entries))
	}
}

func TestDriftedCount_AfterEviction(t *testing.T) {
	w, _ := window.New(50 * time.Millisecond)
	w.Add(makeResults(true)) // should be evicted
	time.Sleep(80 * time.Millisecond)
	w.Add(makeResults(false)) // no drift

	if got := w.DriftedCount(); got != 0 {
		t.Fatalf("expected 0 after eviction, got %d", got)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	w, _ := window.New(time.Minute)
	w.Add(makeResults(false))

	a := w.Entries()
	a[0].Service = "mutated"

	b := w.Entries()
	if b[0].Service == "mutated" {
		t.Fatal("Entries should return an independent copy")
	}
}
