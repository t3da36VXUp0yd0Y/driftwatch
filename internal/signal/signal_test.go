package signal_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/signal"
)

func makeResults(drifted ...bool) []drift.Result {
	names := []string{"alpha", "beta", "gamma"}
	out := make([]drift.Result, 0, len(drifted))
	for i, d := range drifted {
		out = append(out, drift.Result{
			Service: names[i%len(names)],
			Drifted: d,
		})
	}
	return out
}

func TestCompute_NoDrift(t *testing.T) {
	results := makeResults(false, false)
	scores := signal.Compute(results, signal.DefaultWeights())
	for _, s := range scores {
		if s.Drifted {
			t.Errorf("expected no drift for %s", s.Service)
		}
		if s.Value != 0 {
			t.Errorf("expected zero score for %s, got %.2f", s.Service, s.Value)
		}
	}
}

func TestCompute_DriftedServiceHasPositiveScore(t *testing.T) {
	results := []drift.Result{
		{Service: "web", Drifted: true},
		{Service: "db", Drifted: false},
	}
	weights := signal.DefaultWeights()
	scores := signal.Compute(results, weights)

	svcMap := make(map[string]signal.Score)
	for _, s := range scores {
		svcMap[s.Service] = s
	}

	if !svcMap["web"].Drifted {
		t.Error("expected web to be marked drifted")
	}
	if svcMap["web"].Value <= 0 {
		t.Error("expected positive score for drifted service")
	}
	if svcMap["db"].Drifted {
		t.Error("expected db to be healthy")
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	scores := signal.Compute(nil, signal.DefaultWeights())
	if len(scores) != 0 {
		t.Errorf("expected empty scores, got %d", len(scores))
	}
}

func TestCompute_UnknownServiceUsesDefaultFactor(t *testing.T) {
	results := []drift.Result{
		{Service: "mystery", Drifted: true},
	}
	scores := signal.Compute(results, signal.DefaultWeights())
	if len(scores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(scores))
	}
	if scores[0].Value != 0.5 {
		t.Errorf("expected default factor 0.5, got %.2f", scores[0].Value)
	}
}

func TestWrite_ContainsServiceName(t *testing.T) {
	scores := []signal.Score{
		{Service: "api", Value: 1.4, Drifted: true},
		{Service: "cache", Value: 0, Drifted: false},
	}
	var buf bytes.Buffer
	signal.Write(&buf, scores)
	out := buf.String()

	if !strings.Contains(out, "api") {
		t.Error("expected output to contain 'api'")
	}
	if !strings.Contains(out, "drifted") {
		t.Error("expected output to contain 'drifted'")
	}
	if !strings.Contains(out, "ok") {
		t.Error("expected output to contain 'ok'")
	}
}
