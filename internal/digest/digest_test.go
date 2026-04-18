package digest_test

import (
	"testing"

	"github.com/driftwatch/internal/digest"
	"github.com/driftwatch/internal/drift"
)

func makeResults(services []string, drifted bool) []drift.Result {
	out := make([]drift.Result, len(services))
	for i, s := range services {
		out[i] = drift.Result{Service: s, Drifted: drifted, Reason: "image mismatch"}
	}
	return out
}

func TestCompute_Deterministic(t *testing.T) {
	results := makeResults([]string{"api", "worker"}, true)
	a := digest.Compute(results)
	b := digest.Compute(results)
	if a != b {
		t.Fatalf("expected identical digests, got %s and %s", a, b)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := digest.Compute([]drift.Result{
		{Service: "api", Drifted: true, Reason: "image mismatch"},
		{Service: "worker", Drifted: false, Reason: ""},
	})
	b := digest.Compute([]drift.Result{
		{Service: "worker", Drifted: false, Reason: ""},
		{Service: "api", Drifted: true, Reason: "image mismatch"},
	})
	if a != b {
		t.Fatalf("expected order-independent digests, got %s and %s", a, b)
	}
}

func TestCompute_DifferentResults(t *testing.T) {
	a := digest.Compute(makeResults([]string{"api"}, true))
	b := digest.Compute(makeResults([]string{"api"}, false))
	if a == b {
		t.Fatal("expected different digests for different drift states")
	}
}

func TestChanged_True(t *testing.T) {
	results := makeResults([]string{"api"}, true)
	prev := digest.Compute(makeResults([]string{"api"}, false))
	if !digest.Changed(prev, results) {
		t.Fatal("expected Changed to return true")
	}
}

func TestChanged_False(t *testing.T) {
	results := makeResults([]string{"api"}, true)
	prev := digest.Compute(results)
	if digest.Changed(prev, results) {
		t.Fatal("expected Changed to return false")
	}
}

func TestCompute_Empty(t *testing.T) {
	a := digest.Compute(nil)
	b := digest.Compute([]drift.Result{})
	if a != b {
		t.Fatalf("expected same digest for empty inputs, got %s and %s", a, b)
	}
}
