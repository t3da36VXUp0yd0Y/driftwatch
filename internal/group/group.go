package group

import (
	"sort"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Group holds a named collection of drift results.
type Group struct {
	Name    string
	Results []drift.Result
}

// DriftedCount returns the number of drifted results in the group.
func (g Group) DriftedCount() int {
	count := 0
	for _, r := range g.Results {
		if r.Drifted {
			count++
		}
	}
	return count
}

// ByLabel groups results by the value of a given label key.
// Results without the label are placed in a "__unset__" group.
func ByLabel(results []drift.Result, key string) []Group {
	m := make(map[string][]drift.Result)
	for _, r := range results {
		val, ok := r.Labels[key]
		if !ok || strings.TrimSpace(val) == "" {
			val = "__unset__"
		}
		m[val] = append(m[val], r)
	}

	groups := make([]Group, 0, len(m))
	for name, res := range m {
		groups = append(groups, Group{Name: name, Results: res})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}
