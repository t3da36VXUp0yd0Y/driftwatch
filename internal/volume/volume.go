package volume

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Diff represents a single volume mount discrepancy.
type Diff struct {
	Mount    string
	Expected string
	Actual   string
}

func (d Diff) String() string {
	if d.Actual == "" {
		return fmt.Sprintf("mount %q: expected %q, not present", d.Mount, d.Expected)
	}
	return fmt.Sprintf("mount %q: expected %q, got %q", d.Mount, d.Expected, d.Actual)
}

// Check compares declared volume mounts against those found in drift results.
// expected maps service name -> slice of "host:container" mount strings.
func Check(results []drift.Result, expected map[string][]string) []Diff {
	var diffs []Diff
	for _, r := range results {
		declared, ok := expected[r.Service]
		if !ok {
			continue
		}
		actualSet := toSet(r.Mounts)
		for _, m := range declared {
			if _, found := actualSet[m]; !found {
				diffs = append(diffs, Diff{
					Mount:    mountKey(m),
					Expected: m,
					Actual:   "",
				})
			}
		}
	}
	return diffs
}

// Write prints volume diffs to w.
func Write(w io.Writer, diffs []Diff) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "volumes: no drift detected")
		return
	}
	for _, d := range diffs {
		fmt.Fprintln(w, d.String())
	}
}

func toSet(mounts []string) map[string]struct{} {
	s := make(map[string]struct{}, len(mounts))
	for _, m := range mounts {
		s[m] = struct{}{}
	}
	return s
}

func mountKey(m string) string {
	parts := strings.SplitN(m, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return m
}
