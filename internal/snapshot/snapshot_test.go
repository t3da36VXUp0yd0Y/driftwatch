package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/snapshot"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{
			Service:        "api",
			ExpectedImage:  "api:v1",
			ActualImage:    "api:v1",
			Drifted:        false,
		},
		{
			Service:        "worker",
			ExpectedImage:  "worker:v2",
			ActualImage:    "worker:v1",
			Drifted:        true,
		},
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New(makeResults())
	after := time.Now().UTC()

	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("expected CapturedAt between %v and %v, got %v", before, after, s.CapturedAt)
	}
}

func TestNew_StoresResults(t *testing.T) {
	results := makeResults()
	s := snapshot.New(results)

	if len(s.Results) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(s.Results))
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	s := snapshot.New(makeResults())

	path, err := snapshot.Save(dir, s)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", path)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "snapshots")
	s := snapshot.New(makeResults())

	_, err := snapshot.Save(dir, s)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory to be created at %s", dir)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	original := snapshot.New(makeResults())
	original.Meta["env"] = "staging"

	path, err := snapshot.Save(dir, original)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(loaded.Results) != len(original.Results) {
		t.Errorf("expected %d results, got %d", len(original.Results), len(loaded.Results))
	}
	if loaded.Meta["env"] != "staging" {
		t.Errorf("expected meta env=staging, got %q", loaded.Meta["env"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snapshot.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
