package depgraph_test

import (
	"sort"
	"testing"

	"github.com/driftwatch/internal/depgraph"
	"github.com/driftwatch/internal/drift"
)

func makeResults(drifted ...string) []drift.Result {
	all := []string{"api", "worker", "db", "cache"}
	set := make(map[string]struct{})
	for _, s := range drifted {
		set[s] = struct{}{}
	}
	var results []drift.Result
	for _, s := range all {
		_, d := set[s]
		results = append(results, drift.Result{Service: s, Drifted: d})
	}
	return results
}

func sorted(s []string) []string { sort.Strings(s); return s }

func TestAffected_NoDrift(t *testing.T) {
	g := depgraph.New()
	g.AddDependency("api", "db")
	got := g.Affected(makeResults())
	if len(got) != 0 {
		t.Fatalf("expected no affected services, got %v", got)
	}
}

func TestAffected_DirectDependency(t *testing.T) {
	g := depgraph.New()
	g.AddDependency("api", "db")
	g.AddDependency("worker", "cache")

	got := sorted(g.Affected(makeResults("db")))
	if len(got) != 1 || got[0] != "api" {
		t.Fatalf("expected [api], got %v", got)
	}
}

func TestAffected_MultipleDependents(t *testing.T) {
	g := depgraph.New()
	g.AddDependency("api", "db")
	g.AddDependency("worker", "db")

	got := sorted(g.Affected(makeResults("db")))
	if len(got) != 2 || got[0] != "api" || got[1] != "worker" {
		t.Fatalf("expected [api worker], got %v", got)
	}
}

func TestAffected_NoDependencyForDriftedService(t *testing.T) {
	g := depgraph.New()
	g.AddDependency("api", "cache")

	got := g.Affected(makeResults("db"))
	if len(got) != 0 {
		t.Fatalf("expected no affected services, got %v", got)
	}
}

func TestDependencies_Known(t *testing.T) {
	g := depgraph.New()
	g.AddDependency("api", "db")
	g.AddDependency("api", "cache")

	got := sorted(g.Dependencies("api"))
	if len(got) != 2 || got[0] != "cache" || got[1] != "db" {
		t.Fatalf("expected [cache db], got %v", got)
	}
}

func TestDependencies_Unknown(t *testing.T) {
	g := depgraph.New()
	if deps := g.Dependencies("unknown"); deps != nil {
		t.Fatalf("expected nil, got %v", deps)
	}
}
