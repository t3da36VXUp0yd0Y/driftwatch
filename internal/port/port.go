package port

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Diff represents a port binding mismatch for a single service.
type Diff struct {
	Service  string
	Expected []string
	Actual   []string
}

// String returns a human-readable description of the port diff.
func (d Diff) String() string {
	return fmt.Sprintf(
		"service %q: expected ports %v, got %v",
		d.Service, d.Expected, d.Actual,
	)
}

// Check compares expected port bindings from config against those reported
// by the running container. It returns a Diff for every service whose ports
// do not match and a boolean indicating whether any drift was found.
func Check(results []drift.Result, expected map[string][]string) ([]Diff, bool) {
	var diffs []Diff

	for _, r := range results {
		want, ok := expected[r.Service]
		if !ok {
			continue
		}
		if !portsEqual(want, r.Ports) {
			diffs = append(diffs, Diff{
				Service:  r.Service,
				Expected: want,
				Actual:   r.Ports,
			})
		}
	}

	return diffs, len(diffs) > 0
}

// Write prints each Diff to the provided writer.
func Write(w interface{ WriteString(string) (int, error) }, diffs []Diff) error {
	if len(diffs) == 0 {
		_, err := w.WriteString("ports: no drift detected\n")
		return err
	}
	var sb strings.Builder
	for _, d := range diffs {
		sb.WriteString("PORT DRIFT: ")
		sb.WriteString(d.String())
		sb.WriteByte('\n')
	}
	_, err := w.WriteString(sb.String())
	return err
}

func portsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[string]struct{}, len(a))
	for _, p := range a {
		set[p] = struct{}{}
	}
	for _, p := range b {
		if _, ok := set[p]; !ok {
			return false
		}
	}
	return true
}
