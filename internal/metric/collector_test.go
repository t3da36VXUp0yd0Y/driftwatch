package metric

import (
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func makeResults(drifted int, healthy int) []drift.Result {
	var results []drift.Result
	for i := 0; i < drifted; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < healthy; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: false})
	}
	return results
}

func TestRecordResults_Counts(t *testing.T) {
	c := NewCollector()
	c.RecordResults(makeResults(3, 2))

	if got := c.Get(KeyChecksTotal); got != 5 {
		t.Errorf("checks_total: want 5, got %d", got)
	}
	if got := c.Get(KeyDriftTotal); got != 3 {
		t.Errorf("drift_total: want 3, got %d", got)
	}
	if got := c.Get(KeyHealthy); got != 2 {
		t.Errorf("healthy_total: want 2, got %d", got)
	}
}

func TestRecordError_Increments(t *testing.T) {
	c := NewCollector()
	c.RecordError()
	c.RecordError()
	if got := c.Get(KeyErrors); got != 2 {
		t.Errorf("errors_total: want 2, got %d", got)
	}
}

func TestRecordResults_Empty(t *testing.T) {
	c := NewCollector()
	c.RecordResults(nil)
	if got := c.Get(KeyChecksTotal); got != 0 {
		t.Errorf("expected 0 checks, got %d", got)
	}
}
