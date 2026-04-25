package dedupe_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/dedupe"
	"github.com/example/driftwatch/internal/drift"
)

func makeResults(services []string, drifted bool) []drift.Result {
	out := make([]drift.Result, len(services))
	for i, s := range services {
		out[i] = drift.Result{Service: s, Drifted: drifted}
	}
	return out
}

func TestNew_InvalidTTL(t *testing.T) {
	_, err := dedupe.New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	_, err = dedupe.New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestNew_Valid(t *testing.T) {
	d, err := dedupe.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestFilter_NoDrift_AlwaysPassThrough(t *testing.T) {
	d, _ := dedupe.New(time.Minute)
	results := makeResults([]string{"web", "api"}, false)
	out := d.Filter(results)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestFilter_DriftedFirstTimePasses(t *testing.T) {
	d, _ := dedupe.New(time.Minute)
	results := makeResults([]string{"web"}, true)
	out := d.Filter(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestFilter_DriftedSuppressedWithinTTL(t *testing.T) {
	d, _ := dedupe.New(time.Hour)
	results := makeResults([]string{"web"}, true)
	d.Filter(results) // first pass — records timestamp
	out := d.Filter(results) // second pass — should be suppressed
	if len(out) != 0 {
		t.Fatalf("expected 0 results (suppressed), got %d", len(out))
	}
}

func TestFilter_DriftedAllowedAfterTTL(t *testing.T) {
	d, _ := dedupe.New(50 * time.Millisecond)
	results := makeResults([]string{"web"}, true)
	d.Filter(results)
	time.Sleep(60 * time.Millisecond)
	out := d.Filter(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result after TTL expiry, got %d", len(out))
	}
}

func TestFlush_ResetsState(t *testing.T) {
	d, _ := dedupe.New(time.Hour)
	results := makeResults([]string{"web"}, true)
	d.Filter(results)
	d.Flush()
	out := d.Filter(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result after flush, got %d", len(out))
	}
}

func TestEntries_ReturnsTracked(t *testing.T) {
	d, _ := dedupe.New(time.Hour)
	results := makeResults([]string{"web", "db"}, true)
	d.Filter(results)
	entries := d.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
