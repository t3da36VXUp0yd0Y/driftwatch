package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func makeResultsForSummary() []drift.DetectResult {
	return []drift.DetectResult{
		{Service: "web", Running: true, ImageDrift: false},
		{Service: "api", Running: true, ImageDrift: true},
		{Service: "worker", Running: false, ImageDrift: false},
		{Service: "cache", Running: true, ImageDrift: false},
	}
}

func TestNewSummary_Counts(t *testing.T) {
	results := makeResultsForSummary()
	s := NewSummary(results)

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Healthy != 2 {
		t.Errorf("expected Healthy=2, got %d", s.Healthy)
	}
	if s.Drifted != 2 {
		t.Errorf("expected Drifted=2, got %d", s.Drifted)
	}
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
}

func TestNewSummary_AllHealthy(t *testing.T) {
	results := []drift.DetectResult{
		{Service: "a", Running: true, ImageDrift: false},
		{Service: "b", Running: true, ImageDrift: false},
	}
	s := NewSummary(results)
	if s.Drifted != 0 || s.Missing != 0 || s.Healthy != 2 {
		t.Errorf("unexpected summary for all-healthy results: %+v", s)
	}
}

func TestWriteSummary_Output(t *testing.T) {
	s := Summary{Total: 3, Healthy: 2, Drifted: 1, Missing: 0}
	var buf bytes.Buffer
	if err := WriteSummary(&buf, s); err != nil {
		t.Fatalf("WriteSummary returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Summary", "Total services", "Healthy", "Drifted", "Missing"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing expected token %q\nGot:\n%s", want, out)
		}
	}
}
