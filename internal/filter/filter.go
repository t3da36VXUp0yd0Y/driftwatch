// Package filter provides functionality to filter drift results
// based on user-defined criteria such as service name patterns or drift status.
package filter

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Options holds the filtering criteria applied to drift results.
type Options struct {
	// OnlyDrifted restricts results to services with detected drift.
	OnlyDrifted bool
	// ServiceNames filters results to only the named services (case-insensitive).
	// An empty slice means no name filtering is applied.
	ServiceNames []string
}

// Apply filters the given slice of drift results according to the provided
// Options and returns a new slice containing only the matching results.
func Apply(results []drift.Result, opts Options) []drift.Result {
	filtered := make([]drift.Result, 0, len(results))

	for _, r := range results {
		if opts.OnlyDrifted && !r.Drifted {
			continue
		}
		if len(opts.ServiceNames) > 0 && !matchesName(r.Service, opts.ServiceNames) {
			continue
		}
		filtered = append(filtered, r)
	}

	return filtered
}

// matchesName reports whether service matches any of the provided names
// using a case-insensitive comparison.
func matchesName(service string, names []string) bool {
	lower := strings.ToLower(service)
	for _, n := range names {
		if strings.ToLower(n) == lower {
			return true
		}
	}
	return false
}
