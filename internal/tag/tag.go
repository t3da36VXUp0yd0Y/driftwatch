// Package tag provides utilities for matching drift results against
// user-defined label or tag filters.
package tag

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Filter holds a set of required key=value tag pairs.
type Filter struct {
	tags map[string]string
}

// New creates a Filter from a slice of "key=value" strings.
// Entries that do not contain "=" are silently ignored.
func New(pairs []string) *Filter {
	tags := make(map[string]string, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		tags[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return &Filter{tags: tags}
}

// Match reports whether a drift.Result carries all tags in the filter.
// An empty filter matches every result.
func (f *Filter) Match(r drift.Result) bool {
	for k, want := range f.tags {
		got, ok := r.Labels[k]
		if !ok || !strings.EqualFold(got, want) {
			return false
		}
	}
	return true
}

// Apply returns only those results that satisfy the filter.
func (f *Filter) Apply(results []drift.Result) []drift.Result {
	if len(f.tags) == 0 {
		return results
	}
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if f.Match(r) {
			out = append(out, r)
		}
	}
	return out
}
