package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/baseline"
	"github.com/user/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", ExpectedImage: "api:1.0", RunningImage: "api:1.0", Drifted: false},
		{Service: "worker", ExpectedImage: "worker:2.0", RunningImage: "worker:1.9", Drifted: true},
		{Service: "cache", ExpectedImage: "redis:7", RunningImage: "redis:7", Drifted: false},
	}
}

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestCapture_OnlyHealthy(t *testing.T) {
	b := baseline.Capture(makeResults())
	if len(b.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(b.Entries))
	}
	for _, e := range b.Entries {
		if e.Service == "worker" {
			t.Errorf("drifted service 'worker' should not be in baseline")
		}
	}
}

func TestCapture_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	b := baseline.Capture(makeResults())
	if b.CapturedAt.Before(before) {
		t.Errorf("CapturedAt should be set to a recent time")
	}
}

func TestSave_And_Load(t *testing.T) {
	path := tempPath(t)
	b := baseline.Capture(makeResults())

	if err := baseline.Save(b, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Entries) != len(b.Entries) {
		t.Errorf("expected %d entries, got %d", len(b.Entries), len(loaded.Entries))
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	b, err := baseline.Load("/nonexistent/path/baseline.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
	if len(b.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(b.Entries))
	}
}

func TestSave_CreatesNestedDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "dir", "baseline.json")
	b := baseline.Capture(makeResults())
	if err := baseline.Save(b, path); err != nil {
		t.Fatalf("Save with nested path failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at nested path: %v", err)
	}
}
