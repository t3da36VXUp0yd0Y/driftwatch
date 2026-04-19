package severity

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity of detected drift.
type Level int

const (
	LevelNone Level = iota
	LevelLow
	LevelMedium
	LevelHigh
)

var levelNames = map[Level]string{
	LevelNone:   "none",
	LevelLow:    "low",
	LevelMedium: "medium",
	LevelHigh:   "high",
}

func (l Level) String() string {
	if s, ok := levelNames[l]; ok {
		return s
	}
	return fmt.Sprintf("Level(%d)", int(l))
}

// Thresholds defines drift count boundaries for severity classification.
type Thresholds struct {
	Low    int
	Medium int
	High   int
}

// DefaultThresholds returns sensible default thresholds.
func DefaultThresholds() Thresholds {
	return Thresholds{Low: 1, Medium: 3, High: 6}
}

// Evaluate returns the severity Level based on the number of drifted
// services in results and the provided thresholds.
func Evaluate(results []drift.Result, t Thresholds) Level {
	count := 0
	for _, r := range results {
		if r.Drifted {
			count++
		}
	}
	switch {
	case t.High > 0 && count >= t.High:
		return LevelHigh
	case t.Medium > 0 && count >= t.Medium:
		return LevelMedium
	case t.Low > 0 && count >= t.Low:
		return LevelLow
	default:
		return LevelNone
	}
}
