package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/audit"
	"github.com/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "web", Drifted: false, Expected: "nginx:1.25", Actual: "nginx:1.25"},
		{Service: "api", Drifted: true, Expected: "go-api:2.0", Actual: "go-api:1.9"},
	}
}

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := audit.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := audit.New("/tmp/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestRecord_And_Load(t *testing.T) {
	path := tempPath(t)
	l, _ := audit.New(path)

	if err := l.Record(makeResults()); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].Service != "api" {
		t.Errorf("expected service api, got %s", entries[1].Service)
	}
	if !entries[1].Drifted {
		t.Error("expected drifted=true for api")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	entries, err := audit.Load("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestRecord_Appends(t *testing.T) {
	path := tempPath(t)
	l, _ := audit.New(path)

	results := makeResults()
	_ = l.Record(results[:1])
	_ = l.Record(results[1:])

	entries, _ := audit.Load(path)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after two appends, got %d", len(entries))
	}
}

func TestRecord_EmptyResults(t *testing.T) {
	path := tempPath(t)
	l, _ := audit.New(path)

	if err := l.Record(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Log("file not created for empty results — acceptable")
	}
}
