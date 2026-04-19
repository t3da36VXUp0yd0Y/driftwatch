// Package depgraph builds a dependency graph between services and detects
// which services are affected when a drifted service is a dependency of others.
package depgraph

import "github.com/driftwatch/internal/drift"

// Graph holds directed edges: dependent -> set of dependencies.
type Graph struct {
	edges map[string]map[string]struct{}
}

// New returns an empty Graph.
func New() *Graph {
	return &Graph{edges: make(map[string]map[string]struct{})}
}

// AddDependency records that `dependent` depends on `dependency`.
func (g *Graph) AddDependency(dependent, dependency string) {
	if _, ok := g.edges[dependent]; !ok {
		g.edges[dependent] = make(map[string]struct{})
	}
	g.edges[dependent][dependency] = struct{}{}
}

// Affected returns all service names that directly or transitively depend on
// any drifted service found in results.
func (g *Graph) Affected(results []drift.Result) []string {
	drifted := make(map[string]struct{})
	for _, r := range results {
		if r.Drifted {
			drifted[r.Service] = struct{}{}
		}
	}

	seen := make(map[string]struct{})
	for dependent, deps := range g.edges {
		for dep := range deps {
			if _, ok := drifted[dep]; ok {
				seen[dependent] = struct{}{}
			}
		}
	}

	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	return out
}

// Dependencies returns the direct dependencies of the given service.
func (g *Graph) Dependencies(service string) []string {
	deps, ok := g.edges[service]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(deps))
	for d := range deps {
		out = append(out, d)
	}
	return out
}
