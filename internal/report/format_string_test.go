package report_test

import (
	"testing"

	"github.com/driftwatch/internal/report"
)

func TestParseFormat_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected report.Format
	}{
		{"text", report.FormatText},
		{"json", report.FormatJSON},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := report.ParseFormat(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := report.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestFormat_String(t *testing.T) {
	if report.FormatText.String() != "text" {
		t.Errorf("expected 'text', got %q", report.FormatText.String())
	}
	if report.FormatJSON.String() != "json" {
		t.Errorf("expected 'json', got %q", report.FormatJSON.String())
	}
}
