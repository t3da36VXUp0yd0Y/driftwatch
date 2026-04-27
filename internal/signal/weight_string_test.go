package signal_test

import (
	"testing"

	"github.com/driftwatch/internal/signal"
)

func TestWeight_String(t *testing.T) {
	w := signal.Weight{Name: "image", Factor: 1.0}
	got := w.String()
	if got != "image(1.00)" {
		t.Errorf("unexpected string: %q", got)
	}
}

func TestParseWeight_Valid(t *testing.T) {
	cases := []struct {
		input  string
		name   string
		factor float64
	}{
		{"image=1.0", "image", 1.0},
		{"env=0.8", "env", 0.8},
		{"port=0", "port", 0},
	}
	for _, tc := range cases {
		w, err := signal.ParseWeight(tc.input)
		if err != nil {
			t.Errorf("ParseWeight(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if w.Name != tc.name {
			t.Errorf("expected name %q, got %q", tc.name, w.Name)
		}
		if w.Factor != tc.factor {
			t.Errorf("expected factor %.2f, got %.2f", tc.factor, w.Factor)
		}
	}
}

func TestParseWeight_Invalid(t *testing.T) {
	cases := []string{
		"",
		"noequals",
		"=1.0",
		"name=",
		"name=-1",
		"name=abc",
	}
	for _, tc := range cases {
		_, err := signal.ParseWeight(tc)
		if err == nil {
			t.Errorf("ParseWeight(%q) expected error, got nil", tc)
		}
	}
}
