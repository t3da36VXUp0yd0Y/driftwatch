package diff_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/diff"
)

func TestField_String(t *testing.T) {
	f := diff.Field{Name: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}
	got := f.String()
	if !strings.Contains(got, "image") || !strings.Contains(got, "nginx:1.25") {
		t.Errorf("unexpected field string: %s", got)
	}
}

func TestCompare_NoDiff(t *testing.T) {
	r := diff.Compare("web", map[string][2]string{
		"image": {"nginx:1.25", "nginx:1.25"},
	})
	if r.HasDiff() {
		t.Error("expected no diff")
	}
	if !strings.Contains(r.Summary(), "no differences") {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_WithDiff(t *testing.T) {
	r := diff.Compare("api", map[string][2]string{
		"image":  {"myapp:2", "myapp:1"},
		"status": {"running", "running"},
	})
	if !r.HasDiff() {
		t.Error("expected diff")
	}
	if len(r.Fields) != 1 {
		t.Errorf("expected 1 diffed field, got %d", len(r.Fields))
	}
	if r.Fields[0].Name != "image" {
		t.Errorf("expected field name 'image', got %s", r.Fields[0].Name)
	}
}

func TestResult_Summary_WithDiff(t *testing.T) {
	r := diff.Compare("svc", map[string][2]string{
		"image": {"a", "b"},
		"tag":   {"x", "y"},
	})
	if !strings.Contains(r.Summary(), "2 field(s)") {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_EmptyPairs(t *testing.T) {
	r := diff.Compare("empty", map[string][2]string{})
	if r.HasDiff() {
		t.Error("expected no diff for empty pairs")
	}
}
