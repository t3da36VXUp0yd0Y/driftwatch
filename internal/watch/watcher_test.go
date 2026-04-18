package watch

import (
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNew_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, "services: []")
	_, err := New(path, 0, func() {})
	if err != ErrInvalidInterval {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestNew_NilCallback(t *testing.T) {
	path := writeTempConfig(t, "services: []")
	_, err := New(path, time.Second, nil)
	if err != ErrNilCallback {
		t.Fatalf("expected ErrNilCallback, got %v", err)
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New("/nonexistent/path.yaml", time.Second, func() {})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_Valid(t *testing.T) {
	path := writeTempConfig(t, "services: []")
	w, err := New(path, time.Millisecond*50, func() {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestStart_DetectsChange(t *testing.T) {
	path := writeTempConfig(t, "services: []")
	var called atomic.Int32
	w, err := New(path, time.Millisecond*30, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	go w.Start()
	defer w.Stop()

	time.Sleep(time.Millisecond * 20)
	if err := os.WriteFile(path, []byte("services: [updated]"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	time.Sleep(time.Millisecond * 100)
	if called.Load() == 0 {
		t.Error("expected onChange to be called after file change")
	}
}

func TestStart_NoChangeNoCallback(t *testing.T) {
	path := writeTempConfig(t, "services: []")
	var called atomic.Int32
	w, err := New(path, time.Millisecond*30, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	go w.Start()
	time.Sleep(time.Millisecond * 100)
	w.Stop()
	if called.Load() != 0 {
		t.Errorf("expected no calls, got %d", called.Load())
	}
}
