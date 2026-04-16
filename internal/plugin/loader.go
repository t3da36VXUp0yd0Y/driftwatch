package plugin

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// ExecPlugin wraps an external binary as a CheckFn.
// The binary must print one JSON drift.Result per line to stdout.
// This is a lightweight subprocess-based extension point.
func ExecPlugin(binaryPath string) CheckFn {
	return func() ([]drift.Result, error) {
		out, err := exec.Command(binaryPath).Output() //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("exec plugin %q: %w", binaryPath, err)
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var results []drift.Result
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var r drift.Result
			if err := parseResultLine(line, &r); err != nil {
				return nil, fmt.Errorf("exec plugin %q: parse: %w", binaryPath, err)
			}
			results = append(results, r)
		}
		return results, nil
	}
}

// parseResultLine decodes a JSON line into a drift.Result.
func parseResultLine(line string, r *drift.Result) error {
	import_json := func() error {
		// inline to avoid import cycle; use encoding/json via a local import.
		return decodeJSON([]byte(line), r)
	}
	return import_json()
}
