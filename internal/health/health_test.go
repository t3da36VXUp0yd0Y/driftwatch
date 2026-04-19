package health_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/health"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Running: true, Drifted: false},
		{Service: "worker", Running: true, Drifted: true, Reason: "image mismatch"},
		{Service: "db", Running: false, Drifted: false},
	}
}

func TestCheck_HealthyService(t *testing.T) {
	results := health.Check(makeResults())
	if results[0].Status != health.StatusHealthy {
		t.Errorf("expected healthy, got %s", results[0].Status)
	}
}

func TestCheck_DriftedIsUnhealthy(t *testing.T) {
	results := health.Check(makeResults())
	if results[1].Status != health.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", results[1].Status)
	}
	if results[1].Detail == "" {
		t.Error("expected detail to be set for drifted service")
	}
}

func TestCheck_NotRunningIsUnhealthy(t *testing.T) {
	results := health.Check(makeResults())
	if results[2].Status != health.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", results[2].Status)
	}
}

func TestCheck_EmptyResults(t *testing.T) {
	results := health.Check(nil)
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %d", len(results))
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    health.Status
		want string
	}{
		{health.StatusHealthy, "healthy"},
		{health.StatusUnhealthy, "unhealthy"},
		{health.StatusStarting, "starting"},
		{health.StatusUnknown, "unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}

func TestWrite_Output(t *testing.T) {
	results := health.Check(makeResults())
	var buf bytes.Buffer
	health.Write(&buf, results)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output from Write")
	}
	if !bytes.Contains([]byte(out), []byte("api")) {
		t.Error("expected output to contain service name 'api'")
	}
}
