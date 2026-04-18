package env

import (
	"fmt"
	"io"
)

// Result holds the env drift outcome for a single service.
type Result struct {
	Service string
	Diffs   []Diff
}

// HasDrift returns true when at least one env var differs.
func (r Result) HasDrift() bool {
	return len(r.Diffs) > 0
}

// Check evaluates environment drift for a named service given its declared
// and actual environment maps.
func Check(service string, declared, actual map[string]string) Result {
	return Result{
		Service: service,
		Diffs:   Compare(declared, actual),
	}
}

// Write renders a human-readable env drift report to w.
func Write(w io.Writer, results []Result) {
	for _, r := range results {
		if !r.HasDrift() {
			fmt.Fprintf(w, "[OK]   %s — env vars match\n", r.Service)
			continue
		}
		fmt.Fprintf(w, "[DRIFT] %s — %d env var(s) differ\n", r.Service, len(r.Diffs))
		for _, d := range r.Diffs {
			fmt.Fprintf(w, "        %s\n", d.String())
		}
	}
}
