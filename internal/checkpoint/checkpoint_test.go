package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/checkpoint"
	"github.com/driftwatch/internal/drift"
)

func makeResults(drifted ...string) []drift.Result {
	results := make([]drift.Result, len(drifted))
	for i, s := range drifted {
		results[i] = drift.Result{Service: s, Drifted: true}
	}
	return results
}

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestSave_And_Load(t *testing.T) {
	path := tempPath(t)
	results := makeResults("svc-a", "svc-b")

	if err := checkpoint.Save(path, results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(s.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(s.Results))
	}
	if s.RecordedAt.IsZero() {
		t.Error("RecordedAt should not be zero")
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	s, err := checkpoint.Load("/tmp/driftwatch_no_such_file_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Results) != 0 {
		t.Errorf("expected empty results, got %d", len(s.Results))
	}
}

func TestSave_CreatesParentDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	path := filepath.Join(dir, "cp.json")
	if err := checkpoint.Save(path, nil); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestChanged_NoDrift(t *testing.T) {
	prev := checkpoint.State{RecordedAt: time.Now(), Results: makeResults()}
	if checkpoint.Changed(prev, makeResults()) {
		t.Error("expected no change")
	}
}

func TestChanged_NewDrift(t *testing.T) {
	prev := checkpoint.State{Results: []drift.Result{{Service: "svc-a", Drifted: false}}}
	curr := makeResults("svc-a")
	if !checkpoint.Changed(prev, curr) {
		t.Error("expected change detected")
	}
}

func TestChanged_DriftResolved(t *testing.T) {
	prev := checkpoint.State{Results: makeResults("svc-a")}
	curr := []drift.Result{{Service: "svc-a", Drifted: false}}
	if !checkpoint.Changed(prev, curr) {
		t.Error("expected change detected when drift resolved")
	}
}

func TestChanged_SameDrift(t *testing.T) {
	prev := checkpoint.State{Results: makeResults("svc-a", "svc-b")}
	curr := makeResults("svc-a", "svc-b")
	if checkpoint.Changed(prev, curr) {
		t.Error("expected no change")
	}
}
