package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestCheck_NoDrift(t *testing.T) {
	r := Check("api", map[string]string{"PORT": "80"}, map[string]string{"PORT": "80"})
	if r.HasDrift() {
		t.Error("expected no drift")
	}
	if r.Service != "api" {
		t.Errorf("expected service api, got %q", r.Service)
	}
}

func TestCheck_WithDrift(t *testing.T) {
	r := Check("worker",
		map[string]string{"QUEUE": "high"},
		map[string]string{"QUEUE": "low"},
	)
	if !r.HasDrift() {
		t.Error("expected drift")
	}
	if len(r.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(r.Diffs))
	}
}

func TestWrite_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	results := []Result{
		{Service: "svc", Diffs: nil},
	}
	Write(&buf, results)
	if !strings.Contains(buf.String(), "[OK]") {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}

func TestWrite_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	results := []Result{
		{
			Service: "svc",
			Diffs: []Diff{{Key: "PORT", Expected: "80", Actual: "90"}},
		},
	}
	Write(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[DRIFT]") {
		t.Errorf("expected DRIFT in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
}

func TestWrite_Empty(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, nil)
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil results")
	}
}
