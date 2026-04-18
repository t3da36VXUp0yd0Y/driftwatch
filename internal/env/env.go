// Package env compares expected environment variables against those
// observed on a running container.
package env

import (
	"fmt"
	"sort"
	"strings"
)

// Diff represents a single environment variable discrepancy.
type Diff struct {
	Key      string
	Expected string
	Actual   string
	Missing  bool
}

func (d Diff) String() string {
	if d.Missing {
		return fmt.Sprintf("%s: expected %q but not set", d.Key, d.Expected)
	}
	return fmt.Sprintf("%s: expected %q got %q", d.Key, d.Expected, d.Actual)
}

// Compare checks declared env vars against the actual set observed on a
// container. Only keys present in declared are checked; extra keys on the
// container are ignored.
func Compare(declared, actual map[string]string) []Diff {
	var diffs []Diff
	keys := make([]string, 0, len(declared))
	for k := range declared {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		want := declared[k]
		got, ok := actual[k]
		if !ok {
			diffs = append(diffs, Diff{Key: k, Expected: want, Missing: true})
			continue
		}
		if got != want {
			diffs = append(diffs, Diff{Key: k, Expected: want, Actual: got})
		}
	}
	return diffs
}

// ParseEnvSlice converts a slice of "KEY=VALUE" strings into a map.
func ParseEnvSlice(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
