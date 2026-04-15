package alert_test

import (
	"bytes"
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func makeResults(drifted, healthy int) []drift.Result {
	var results []drift.Result
	for i := 0; i < drifted; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < healthy; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: false})
	}
	return results
}

func TestEvaluate_NoDrift(t *testing.T) {
	r := makeResults(0, 3)
	a := alert.Evaluate(r, alert.Thresholds{Warn: 1, Critical: 3})
	if a.Level != alert.LevelNone {
		t.Errorf("expected LevelNone, got %s", a.Level)
	}
}

func TestEvaluate_WarnThreshold(t *testing.T) {
	r := makeResults(2, 1)
	a := alert.Evaluate(r, alert.Thresholds{Warn: 2, Critical: 5})
	if a.Level != alert.LevelWarning {
		t.Errorf("expected LevelWarning, got %s", a.Level)
	}
	if a.DriftCount != 2 {
		t.Errorf("expected DriftCount 2, got %d", a.DriftCount)
	}
}

func TestEvaluate_CriticalThreshold(t *testing.T) {
	r := makeResults(4, 0)
	a := alert.Evaluate(r, alert.Thresholds{Warn: 2, Critical: 4})
	if a.Level != alert.LevelCritical {
		t.Errorf("expected LevelCritical, got %s", a.Level)
	}
}

func TestEvaluate_ZeroThresholdsNeverAlert(t *testing.T) {
	r := makeResults(10, 0)
	a := alert.Evaluate(r, alert.Thresholds{Warn: 0, Critical: 0})
	if a.Level != alert.LevelNone {
		t.Errorf("expected LevelNone with zero thresholds, got %s", a.Level)
	}
}

func TestWrite_NoneWritesNothing(t *testing.T) {
	var buf bytes.Buffer
	a := alert.Alert{Level: alert.LevelNone}
	if err := alert.Write(&buf, a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for LevelNone, got %q", buf.String())
	}
}

func TestWrite_Warning(t *testing.T) {
	var buf bytes.Buffer
	a := alert.Alert{
		Level:      alert.LevelWarning,
		DriftCount: 2,
		TotalCount: 5,
		Thresholds: alert.Thresholds{Warn: 2, Critical: 4},
	}
	if err := alert.Write(&buf, a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output for LevelWarning")
	}
}
