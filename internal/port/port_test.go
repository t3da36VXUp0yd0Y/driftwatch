package port_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/port"
)

func makeResults(svc string, ports []string) drift.Result {
	return drift.Result{Service: svc, Ports: ports}
}

func TestCheck_NoDrift(t *testing.T) {
	results := []drift.Result{
		makeResults("web", []string{"80:80", "443:443"}),
	}
	expected := map[string][]string{
		"web": {"443:443", "80:80"},
	}
	diffs, drifted := port.Check(results, expected)
	if drifted {
		t.Errorf("expected no drift, got %v", diffs)
	}
	if len(diffs) != 0 {
		t.Errorf("expected empty diffs, got %d", len(diffs))
	}
}

func TestCheck_PortMismatch(t *testing.T) {
	results := []drift.Result{
		makeResults("api", []string{"8080:8080"}),
	}
	expected := map[string][]string{
		"api": {"9090:9090"},
	}
	diffs, drifted := port.Check(results, expected)
	if !drifted {
		t.Fatal("expected drift to be detected")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Service != "api" {
		t.Errorf("unexpected service name %q", diffs[0].Service)
	}
}

func TestCheck_ServiceNotInExpected(t *testing.T) {
	results := []drift.Result{
		makeResults("worker", []string{"5000:5000"}),
	}
	diffs, drifted := port.Check(results, map[string][]string{})
	if drifted || len(diffs) != 0 {
		t.Error("expected no drift for unconfigured service")
	}
}

func TestDiff_String(t *testing.T) {
	d := port.Diff{
		Service:  "web",
		Expected: []string{"80:80"},
		Actual:   []string{"8080:80"},
	}
	s := d.String()
	if !strings.Contains(s, "web") {
		t.Errorf("expected service name in string, got %q", s)
	}
}

func TestWrite_NoDrift(t *testing.T) {
	var sb strings.Builder
	if err := port.Write(&sb, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no drift") {
		t.Errorf("expected no-drift message, got %q", sb.String())
	}
}

func TestWrite_WithDrift(t *testing.T) {
	var sb strings.Builder
	diffs := []port.Diff{{Service: "db", Expected: []string{"5432:5432"}, Actual: []string{"3306:3306"}}}
	if err := port.Write(&sb, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "PORT DRIFT") {
		t.Errorf("expected PORT DRIFT prefix, got %q", sb.String())
	}
}
