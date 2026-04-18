package volume_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/volume"
)

func makeResults(mounts []string) []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Mounts:  mounts,
		},
	}
}

func TestCheck_NoDrift(t *testing.T) {
	results := makeResults([]string{"/data:/app/data", "/logs:/app/logs"})
	expected := map[string][]string{
		"api": {"/data:/app/data", "/logs:/app/logs"},
	}
	diffs := volume.Check(results, expected)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCheck_MissingMount(t *testing.T) {
	results := makeResults([]string{"/data:/app/data"})
	expected := map[string][]string{
		"api": {"/data:/app/data", "/logs:/app/logs"},
	}
	diffs := volume.Check(results, expected)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Expected != "/logs:/app/logs" {
		t.Errorf("unexpected diff expected value: %s", diffs[0].Expected)
	}
}

func TestCheck_ServiceNotInExpected(t *testing.T) {
	results := makeResults([]string{"/data:/app/data"})
	expected := map[string][]string{
		"other": {"/data:/app/data"},
	}
	diffs := volume.Check(results, expected)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs for unlisted service, got %d", len(diffs))
	}
}

func TestDiff_String_Missing(t *testing.T) {
	d := volume.Diff{Mount: "/app/data", Expected: "/data:/app/data", Actual: ""}
	if !strings.Contains(d.String(), "not present") {
		t.Errorf("expected 'not present' in: %s", d.String())
	}
}

func TestWrite_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	volume.Write(&buf, nil)
	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected no drift message, got: %s", buf.String())
	}
}

func TestWrite_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	diffs := []volume.Diff{
		{Mount: "/app/data", Expected: "/data:/app/data", Actual: ""},
	}
	volume.Write(&buf, diffs)
	if !strings.Contains(buf.String(), "/app/data") {
		t.Errorf("expected mount path in output, got: %s", buf.String())
	}
}
