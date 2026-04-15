package alert_test

import (
	"testing"

	"github.com/driftwatch/internal/alert"
)

func TestParseLevel_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected alert.Level
	}{
		{"none", alert.LevelNone},
		{"NONE", alert.LevelNone},
		{"warn", alert.LevelWarning},
		{"WARNING", alert.LevelWarning},
		{"crit", alert.LevelCritical},
		{"CRITICAL", alert.LevelCritical},
	}
	for _, c := range cases {
		got, err := alert.ParseLevel(c.input)
		if err != nil {
			t.Errorf("ParseLevel(%q) unexpected error: %v", c.input, err)
		}
		if got != c.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := alert.ParseLevel("unknown")
	if err == nil {
		t.Error("expected error for unknown level, got nil")
	}
}

func TestLevel_String(t *testing.T) {
	if alert.LevelNone.String() != "NONE" {
		t.Errorf("expected NONE, got %s", alert.LevelNone.String())
	}
	if alert.LevelWarning.String() != "WARNING" {
		t.Errorf("expected WARNING, got %s", alert.LevelWarning.String())
	}
	if alert.LevelCritical.String() != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", alert.LevelCritical.String())
	}
}
