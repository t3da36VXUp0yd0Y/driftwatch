package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/history"
)

func makeResults(drifted bool) []drift.Result {
	return []drift.Result{
		{
			Service:        "web",
			ExpectedImage:  "nginx:1.25",
			ActualImage:    "nginx:1.24",
			Drifted:        drifted,
			ContainerID:    "abc123",
			ContainerName:  "web",
		},
	}
}

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestAppend_CreatesFile(t *testing.T) {
	path := tempPath(t)
	s := history.New(path)

	if err := s.Append(makeResults(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected history file to be created")
	}
}

func TestLoad_EmptyWhenMissing(t *testing.T) {
	s := history.New(tempPath(t))
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestAppend_AccumulatesEntries(t *testing.T) {
	path := tempPath(t)
	s := history.New(path)

	if err := s.Append(makeResults(false)); err != nil {
		t.Fatalf("first append: %v", err)
	}
	if err := s.Append(makeResults(true)); err != nil {
		t.Fatalf("second append: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	path := tempPath(t)
	s := history.New(path)

	_ = s.Append(makeResults(false))
	time.Sleep(time.Millisecond)
	_ = s.Append(makeResults(true))

	entry, ok := s.Latest()
	if !ok {
		t.Fatal("expected Latest to return an entry")
	}
	if !entry.Results[0].Drifted {
		t.Error("expected latest entry to have drifted=true")
	}
}

func TestLatest_FalseWhenEmpty(t *testing.T) {
	s := history.New(tempPath(t))
	_, ok := s.Latest()
	if ok {
		t.Fatal("expected Latest to return false for empty store")
	}
}
