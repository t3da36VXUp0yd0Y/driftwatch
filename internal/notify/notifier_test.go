package notify_test

import (
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/notify"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "web", Expected: "nginx:1.25", Actual: "nginx:1.25", Drifted: false, Reason: ""},
		{Service: "api", Expected: "myapp:2.0", Actual: "myapp:1.9", Drifted: true, Reason: "image mismatch"},
	}
}

func TestNotify_LevelAll_WritesAll(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelAll)
	count, err := n.Notify(makeResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 notifications, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected OK line for healthy service")
	}
	if !strings.Contains(out, "[DRIFT]") {
		t.Error("expected DRIFT line for drifted service")
	}
}

func TestNotify_LevelDriftOnly_WritesOnlyDrifted(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelDriftOnly)
	count, err := n.Notify(makeResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 notification, got %d", count)
	}
	out := buf.String()
	if strings.Contains(out, "[OK]") {
		t.Error("did not expect OK line in drift-only mode")
	}
	if !strings.Contains(out, "service=api") {
		t.Error("expected api service in output")
	}
}

func TestNotify_EmptyResults(t *testing.T) {
	var buf strings.Builder
	n := notify.New(&buf, notify.LevelAll)
	count, err := n.Notify(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 notifications, got %d", count)
	}
}
