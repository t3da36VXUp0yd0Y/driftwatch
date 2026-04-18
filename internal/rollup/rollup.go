// Package rollup aggregates drift results across multiple services into
// grouped summaries, making it easier to identify systemic issues.
package rollup

import (
	"fmt"
	"sort"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Group holds aggregated drift results for a named group.
type Group struct {
	Name    string
	Total   int
	Drifted int
	Services []string
}

// Healthy returns true when no services in the group have drifted.
func (g Group) Healthy() bool { return g.Drifted == 0 }

// String returns a short human-readable summary of the group.
func (g Group) String() string {
	return fmt.Sprintf("%s: %d/%d drifted", g.Name, g.Drifted, g.Total)
}

// ByPrefix groups results by the prefix of each service name, using sep as
// the delimiter. Services without sep are placed in a "default" group.
func ByPrefix(results []drift.Result, sep string) []Group {
	buckets := map[string]*Group{}

	for _, r := range results {
		key := "default"
		if idx := strings.Index(r.Service, sep); idx > 0 {
			key = r.Service[:idx]
		}
		g, ok := buckets[key]
		if !ok {
			g = &Group{Name: key}
			buckets[key] = g
		}
		g.Total++
		g.Services = append(g.Services, r.Service)
		if r.Drifted {
			g.Drifted++
		}
	}

	groups := make([]Group, 0, len(buckets))
	for _, g := range buckets {
		groups = append(groups, *g)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}
