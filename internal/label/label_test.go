package label_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/label"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Labels: map[string]string{"env": "prod", "team": "platform"}},
		{Service: "worker", Labels: map[string]string{"env": "prod", "team": "data"}},
		{Service: "cache", Labels: map[string]string{"env": "staging"}},
	}
}

func TestNew_IgnoresInvalidPairs(t *testing.T) {
	f := label.New([]string{"noequals", "=emptykey", "env=prod"})
	results := f.Apply(makeResults())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	f := label.New(nil)
	results := f.Apply(makeResults())
	if len(results) != 3 {
		t.Fatalf("expected all 3 results, got %d", len(results))
	}
}

func TestMatch_SingleSelector(t *testing.T) {
	f := label.New([]string{"team=platform"})
	results := f.Apply(makeResults())
	if len(results) != 1 || results[0].Service != "api" {
		t.Fatalf("expected only api, got %v", results)
	}
}

func TestMatch_MultipleSelectors(t *testing.T) {
	f := label.New([]string{"env=prod", "team=data"})
	results := f.Apply(makeResults())
	if len(results) != 1 || results[0].Service != "worker" {
		t.Fatalf("expected only worker, got %v", results)
	}
}

func TestMatch_NoMatchReturnsEmpty(t *testing.T) {
	f := label.New([]string{"env=canary"})
	results := f.Apply(makeResults())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestMatch_LabelNotPresent(t *testing.T) {
	f := label.New([]string{"region=us-east"})
	if f.Match(map[string]string{"env": "prod"}) {
		t.Fatal("expected no match when label is absent")
	}
}
