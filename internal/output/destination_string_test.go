package output_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/output"
)

func TestParseDestination_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  output.Destination
	}{
		{"stdout", output.Stdout},
		{"file", output.File},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := output.ParseDestination(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseDestination(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseDestination_Invalid(t *testing.T) {
	_, err := output.ParseDestination("printer")
	if err == nil {
		t.Fatal("expected error for unknown destination, got nil")
	}
}

func TestDestination_String(t *testing.T) {
	cases := []struct {
		d    output.Destination
		want string
	}{
		{output.Stdout, "stdout"},
		{output.File, "file"},
		{output.Destination(42), "Destination(42)"},
	}
	for _, tc := range cases {
		if got := tc.d.String(); got != tc.want {
			t.Errorf("Destination(%d).String() = %q, want %q", int(tc.d), got, tc.want)
		}
	}
}
