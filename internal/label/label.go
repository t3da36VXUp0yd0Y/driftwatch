// Package label provides utilities for filtering drift results by container labels.
package label

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Filter holds the label selectors used to narrow results.
type Filter struct {
	selectors map[string]string
}

// New creates a Filter from a slice of "key=value" strings.
// Entries that do not contain "=" are ignored.
func New(pairs []string) *Filter {
	sel := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || strings.TrimSpace(k) == "" {
			continue
		}
		sel[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return &Filter{selectors: sel}
}

// Match reports whether the given label map satisfies all selectors.
// An empty filter matches everything.
func (f *Filter) Match(labels map[string]string) bool {
	for k, want := range f.selectors {
		got, ok := labels[k]
		if !ok || got != want {
			return false
		}
	}
	return true
}

// Apply returns only those results whose Labels satisfy the filter.
func (f *Filter) Apply(results []drift.Result) []drift.Result {
	if len(f.selectors) == 0 {
		return results
	}
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if f.Match(r.Labels) {
			out = append(out, r)
		}
	}
	return out
}
