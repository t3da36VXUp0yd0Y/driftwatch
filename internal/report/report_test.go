package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func makeResults(drifted bool) []drift.Result {
	image := "nginx:latest"
	if drifted {
		image = "nginx:1.19"
	}
	return []drift.Result{
		{
			ServiceName:   "web",
			ExpectedImage: "nginx:latest",
			ActualImage:   image,
			Running:       true,
		},
	}
}

func TestReport_HasDrift_False(t *testing.T) {
	r := report.New(makeResults(false))
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_HasDrift_True(t *testing.T) {
	r := report.New(makeResults(true))
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestReport_WriteText(t *testing.T) {
	r := report.New(makeResults(true))
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "web") {
		t.Error("expected service name 'web' in output")
	}
	if !strings.Contains(out, "DRIFT") {
		t.Error("expected DRIFT status in output")
	}
}

func TestReport_WriteJSON(t *testing.T) {
	r := report.New(makeResults(true))
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["has_drift"] != true {
		t.Error("expected has_drift to be true")
	}
	if payload["generated_at"] == "" {
		t.Error("expected generated_at to be set")
	}
}

func TestReport_WriteText_NoDrift(t *testing.T) {
	r := report.New(makeResults(false))
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "DRIFT") {
		t.Error("expected no DRIFT status in output")
	}
}
