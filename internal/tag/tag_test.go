package tag_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/tag"
)

func makeResult(labels map[string]string, drifted bool) drift.Result {
	return drift.Result{
		Service: "svc",
		Drifted: drifted,
		Labels:  labels,
	}
}

func TestNew_IgnoresInvalidPairs(t *testing.T) {
	f := tag.New([]string{"noequals", "env=prod"})
	r := makeResult(map[string]string{"env": "prod"}, false)
	if !f.Match(r) {
		t.Fatal("expected match")
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	f := tag.New(nil)
	r := makeResult(nil, false)
	if !f.Match(r) {
		t.Fatal("empty filter should match everything")
	}
}

func TestMatch_AllTagsPresent(t *testing.T) {
	f := tag.New([]string{"env=prod", "team=platform"})
	r := makeResult(map[string]string{"env": "prod", "team": "platform"}, true)
	if !f.Match(r) {
		t.Fatal("expected match")
	}
}

func TestMatch_MissingTag(t *testing.T) {
	f := tag.New([]string{"env=prod"})
	r := makeResult(map[string]string{"team": "platform"}, false)
	if f.Match(r) {
		t.Fatal("expected no match")
	}
}

func TestMatch_CaseInsensitiveValue(t *testing.T) {
	f := tag.New([]string{"env=PROD"})
	r := makeResult(map[string]string{"env": "prod"}, false)
	if !f.Match(r) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestApply_FiltersResults(t *testing.T) {
	f := tag.New([]string{"env=staging"})
	results := []drift.Result{
		makeResult(map[string]string{"env": "staging"}, true),
		makeResult(map[string]string{"env": "prod"}, false),
		makeResult(map[string]string{"env": "staging"}, false),
	}
	out := f.Apply(results)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_EmptyFilter_ReturnsAll(t *testing.T) {
	f := tag.New([]string{})
	results := []drift.Result{
		makeResult(map[string]string{"env": "prod"}, false),
		makeResult(nil, true),
	}
	out := f.Apply(results)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}
