package severity_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/severity"
)

func makeResults(drifted, healthy int) []drift.Result {
	var results []drift.Result
	for i := 0; i < drifted; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < healthy; i++ {
		results = append(results, drift.Result{Service: "ok", Drifted: false})
	}
	return results
}

func TestEvaluate_NoDrift(t *testing.T) {
	res := makeResults(0, 3)
	got := severity.Evaluate(res, severity.DefaultThresholds())
	if got != severity.LevelNone {
		t.Fatalf("expected none, got %s", got)
	}
}

func TestEvaluate_LowThreshold(t *testing.T) {
	res := makeResults(1, 2)
	got := severity.Evaluate(res, severity.DefaultThresholds())
	if got != severity.LevelLow {
		t.Fatalf("expected low, got %s", got)
	}
}

func TestEvaluate_MediumThreshold(t *testing.T) {
	res := makeResults(3, 0)
	got := severity.Evaluate(res, severity.DefaultThresholds())
	if got != severity.LevelMedium {
		t.Fatalf("expected medium, got %s", got)
	}
}

func TestEvaluate_HighThreshold(t *testing.T) {
	res := makeResults(6, 1)
	got := severity.Evaluate(res, severity.DefaultThresholds())
	if got != severity.LevelHigh {
		t.Fatalf("expected high, got %s", got)
	}
}

func TestEvaluate_ZeroThresholdsNeverAlert(t *testing.T) {
	res := makeResults(10, 0)
	got := severity.Evaluate(res, severity.Thresholds{})
	if got != severity.LevelNone {
		t.Fatalf("expected none with zero thresholds, got %s", got)
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level severity.Level
		want  string
	}{
		{severity.LevelNone, "none"},
		{severity.LevelLow, "low"},
		{severity.LevelMedium, "medium"},
		{severity.LevelHigh, "high"},
	}
	for _, c := range cases {
		if got := c.level.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}
