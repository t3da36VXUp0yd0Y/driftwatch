package exit_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/exit"
)

func TestCode_String_Known(t *testing.T) {
	tests := []struct {
		code exit.Code
		want string
	}{
		{exit.OK, "ok"},
		{exit.DriftDetected, "drift_detected"},
		{exit.ConfigError, "config_error"},
		{exit.ClientError, "client_error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.code.String()
			if got != tt.want {
				t.Errorf("Code(%d).String() = %q; want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestCode_String_Unknown(t *testing.T) {
	c := exit.Code(99)
	if got := c.String(); got != "unknown" {
		t.Errorf("unexpected string for unknown code: %q", got)
	}
}

func TestCode_IntValues(t *testing.T) {
	if int(exit.OK) != 0 {
		t.Errorf("OK should be 0, got %d", exit.OK)
	}
	if int(exit.DriftDetected) != 1 {
		t.Errorf("DriftDetected should be 1, got %d", exit.DriftDetected)
	}
	if int(exit.ConfigError) != 2 {
		t.Errorf("ConfigError should be 2, got %d", exit.ConfigError)
	}
	if int(exit.ClientError) != 3 {
		t.Errorf("ClientError should be 3, got %d", exit.ClientError)
	}
}
