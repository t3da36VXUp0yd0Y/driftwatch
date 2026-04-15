package notify_test

import (
	"testing"

	"github.com/example/driftwatch/internal/notify"
)

func TestParseLevel_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  notify.Level
	}{
		{"all", notify.LevelAll},
		{"ALL", notify.LevelAll},
		{"drift-only", notify.LevelDriftOnly},
		{"driftonly", notify.LevelDriftOnly},
		{"DRIFT-ONLY", notify.LevelDriftOnly},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := notify.ParseLevel(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := notify.ParseLevel("verbose")
	if err == nil {
		t.Error("expected error for unknown level")
	}
}

func TestLevel_String(t *testing.T) {
	if got := notify.LevelAll.String(); got != "all" {
		t.Errorf("LevelAll.String() = %q, want %q", got, "all")
	}
	if got := notify.LevelDriftOnly.String(); got != "drift-only" {
		t.Errorf("LevelDriftOnly.String() = %q, want %q", got, "drift-only")
	}
	unknown := notify.Level(99)
	if got := unknown.String(); got != "Level(99)" {
		t.Errorf("unknown.String() = %q, want %q", got, "Level(99)")
	}
}
