package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func makeResults(names []string, drifted map[string]bool) []drift.Result {
	out := make([]drift.Result, 0, len(names))
	for _, n := range names {
		out = append(out, drift.Result{Service: n, Drifted: drifted[n]})
	}
	return out
}

func TestByPrefix_GroupsCorrectly(t *testing.T) {
	results := makeResults(
		[]string{"team-a-api", "team-a-worker", "team-b-api"},
		map[string]bool{"team-a-api": true},
	)
	groups := rollup.ByPrefix(results, "-")
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "team" {
		t.Errorf("unexpected group name: %s", groups[0].Name)
	}
}

func TestByPrefix_DefaultGroup(t *testing.T) {
	results := makeResults([]string{"standalone"}, nil)
	groups := rollup.ByPrefix(results, "-")
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "default" {
		t.Errorf("expected default group, got %s", groups[0].Name)
	}
}

func TestByPrefix_DriftedCount(t *testing.T) {
	results := makeResults(
		[]string{"svc-one", "svc-two", "svc-three"},
		map[string]bool{"svc-one": true, "svc-three": true},
	)
	groups := rollup.ByPrefix(results, "-")
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	g := groups[0]
	if g.Total != 3 {
		t.Errorf("expected Total=3, got %d", g.Total)
	}
	if g.Drifted != 2 {
		t.Errorf("expected Drifted=2, got %d", g.Drifted)
	}
	if g.Healthy() {
		t.Error("expected group to be unhealthy")
	}
}

func TestByPrefix_EmptyResults(t *testing.T) {
	groups := rollup.ByPrefix(nil, "-")
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestGroup_String(t *testing.T) {
	g := rollup.Group{Name: "team", Total: 4, Drifted: 1}
	s := g.String()
	if s != "team: 1/4 drifted" {
		t.Errorf("unexpected string: %s", s)
	}
}
