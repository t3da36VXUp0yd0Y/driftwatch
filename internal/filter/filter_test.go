package filter_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/filter"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Drifted: true, Reason: "image mismatch"},
		{Service: "worker", Drifted: false},
		{Service: "scheduler", Drifted: true, Reason: "not running"},
	}
}

func TestApply_NoFilter(t *testing.T) {
	results := makeResults()
	out := filter.Apply(results, filter.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_OnlyDrifted(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{OnlyDrifted: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(out))
	}
	for _, r := range out {
		if !r.Drifted {
			t.Errorf("expected drifted=true for service %q", r.Service)
		}
	}
}

func TestApply_ServiceNames(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{ServiceNames: []string{"api", "worker"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_ServiceNames_CaseInsensitive(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{ServiceNames: []string{"API"}})
	if len(out) != 1 || out[0].Service != "api" {
		t.Fatalf("expected 1 result for 'API', got %d", len(out))
	}
}

func TestApply_OnlyDrifted_AndServiceNames(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{
		OnlyDrifted:  true,
		ServiceNames: []string{"api"},
	})
	if len(out) != 1 || out[0].Service != "api" {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestApply_NoMatch(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{ServiceNames: []string{"nonexistent"}})
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}
