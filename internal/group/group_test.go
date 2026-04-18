package group_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/group"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Drifted: true, Labels: map[string]string{"env": "prod"}},
		{Service: "worker", Drifted: false, Labels: map[string]string{"env": "prod"}},
		{Service: "cache", Drifted: true, Labels: map[string]string{"env": "staging"}},
		{Service: "db", Drifted: false, Labels: map[string]string{}},
	}
}

func TestByLabel_GroupsCorrectly(t *testing.T) {
	results := makeResults()
	groups := group.ByLabel(results, "env")

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestByLabel_UnsetFallback(t *testing.T) {
	results := makeResults()
	groups := group.ByLabel(results, "env")

	var unset *group.Group
	for i := range groups {
		if groups[i].Name == "__unset__" {
			unset = &groups[i]
		}
	}
	if unset == nil {
		t.Fatal("expected __unset__ group")
	}
	if len(unset.Results) != 1 {
		t.Errorf("expected 1 result in __unset__, got %d", len(unset.Results))
	}
}

func TestByLabel_DriftedCount(t *testing.T) {
	results := makeResults()
	groups := group.ByLabel(results, "env")

	for _, g := range groups {
		if g.Name == "prod" && g.DriftedCount() != 1 {
			t.Errorf("prod: expected 1 drifted, got %d", g.DriftedCount())
		}
		if g.Name == "staging" && g.DriftedCount() != 1 {
			t.Errorf("staging: expected 1 drifted, got %d", g.DriftedCount())
		}
	}
}

func TestByLabel_EmptyResults(t *testing.T) {
	groups := group.ByLabel([]drift.Result{}, "env")
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestByLabel_SortedByName(t *testing.T) {
	results := makeResults()
	groups := group.ByLabel(results, "env")
	for i := 1; i < len(groups); i++ {
		if groups[i-1].Name > groups[i].Name {
			t.Errorf("groups not sorted: %s > %s", groups[i-1].Name, groups[i].Name)
		}
	}
}
