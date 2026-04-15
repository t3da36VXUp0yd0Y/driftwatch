// Package alert provides threshold-based alerting for drift detection results.
// It evaluates drift results against configurable thresholds and emits alerts
// when the number of drifted services exceeds acceptable limits.
package alert

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity of an alert.
type Level int

const (
	LevelNone    Level = iota // no threshold exceeded
	LevelWarning              // warn threshold exceeded
	LevelCritical             // critical threshold exceeded
)

// Thresholds defines the warning and critical drift count limits.
type Thresholds struct {
	Warn     int
	Critical int
}

// Alert holds the result of a threshold evaluation.
type Alert struct {
	Level       Level
	DriftCount  int
	TotalCount  int
	Thresholds  Thresholds
}

// Evaluate inspects drift results against the given thresholds and returns an Alert.
func Evaluate(results []drift.Result, t Thresholds) Alert {
	total := len(results)
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	level := LevelNone
	switch {
	case t.Critical > 0 && drifted >= t.Critical:
		level = LevelCritical
	case t.Warn > 0 && drifted >= t.Warn:
		level = LevelWarning
	}

	return Alert{
		Level:      level,
		DriftCount: drifted,
		TotalCount: total,
		Thresholds: t,
	}
}

// Write formats the alert and writes it to w. Returns nil if level is None.
func Write(w io.Writer, a Alert) error {
	if a.Level == LevelNone {
		return nil
	}
	_, err := fmt.Fprintf(w, "[%s] %d/%d services drifted (warn: %d, critical: %d)\n",
		a.Level, a.DriftCount, a.TotalCount, a.Thresholds.Warn, a.Thresholds.Critical)
	return err
}
