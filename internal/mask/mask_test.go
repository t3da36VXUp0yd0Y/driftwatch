package mask_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/mask"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24", Drifted: true},
		{Service: "api", Field: "env.SECRET_KEY", Expected: "abc123", Actual: "xyz789", Drifted: true},
		{Service: "db", Field: "image", Expected: "postgres:15", Actual: "postgres:15", Drifted: false},
	}
}

func TestNew_IgnoresEmptyFields(t *testing.T) {
	m := mask.New([]mask.Rule{{Field: "", Replacement: "hidden"}})
	results := makeResults()
	out := m.Apply(results)
	if out[0].Expected != results[0].Expected {
		t.Errorf("expected value unchanged, got %q", out[0].Expected)
	}
}

func TestApply_NoRules(t *testing.T) {
	m := mask.New(nil)
	results := makeResults()
	out := m.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
	if out[1].Expected != "abc123" {
		t.Errorf("expected unmasked value, got %q", out[1].Expected)
	}
}

func TestApply_MasksMatchingField(t *testing.T) {
	m := mask.New([]mask.Rule{{Field: "env.SECRET_KEY", Replacement: "[redacted]"}})
	out := m.Apply(makeResults())
	if out[1].Expected != "[redacted]" {
		t.Errorf("expected [redacted], got %q", out[1].Expected)
	}
	if out[1].Actual != "[redacted]" {
		t.Errorf("expected [redacted], got %q", out[1].Actual)
	}
}

func TestApply_DefaultReplacement(t *testing.T) {
	m := mask.New([]mask.Rule{{Field: "env.SECRET_KEY"}})
	out := m.Apply(makeResults())
	if out[1].Expected != mask.DefaultReplacement {
		t.Errorf("expected %q, got %q", mask.DefaultReplacement, out[1].Expected)
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	m := mask.New([]mask.Rule{{Field: "ENV.SECRET_KEY", Replacement: "***"}})
	out := m.Apply(makeResults())
	if out[1].Actual != "***" {
		t.Errorf("expected masked value, got %q", out[1].Actual)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults()
	m := mask.New([]mask.Rule{{Field: "env.SECRET_KEY", Replacement: "hidden"}})
	m.Apply(results)
	if results[1].Expected != "abc123" {
		t.Errorf("original result mutated, got %q", results[1].Expected)
	}
}
