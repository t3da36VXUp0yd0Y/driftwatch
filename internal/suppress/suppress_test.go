package suppress_test

import (
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/driftwatch/driftwatch/internal/suppress"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Drifted: true, Detail: "image mismatch"},
		{Service: "worker", Drifted: true, Detail: "env var changed"},
		{Service: "db", Drifted: false, Detail: ""},
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	s := suppress.New(nil)
	results := makeResults()
	out := s.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_SuppressServiceAllDrift(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "api"},
	}
	s := suppress.New(rules)
	out := s.Apply(makeResults())
	for _, r := range out {
		if r.Service == "api" {
			t.Errorf("expected api to be suppressed")
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results after suppression, got %d", len(out))
	}
}

func TestApply_SuppressMatchingReason(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "worker", Reason: "env var"},
	}
	s := suppress.New(rules)
	out := s.Apply(makeResults())
	for _, r := range out {
		if r.Service == "worker" {
			t.Errorf("expected worker to be suppressed by reason match")
		}
	}
}

func TestApply_ReasonMismatch_NotSuppressed(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "api", Reason: "port mismatch"},
	}
	s := suppress.New(rules)
	out := s.Apply(makeResults())
	found := false
	for _, r := range out {
		if r.Service == "api" {
			found = true
		}
	}
	if !found {
		t.Error("api should not be suppressed when reason does not match")
	}
}

func TestApply_ExpiredRule_NotSuppressed(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "api", ExpiresAt: time.Now().Add(-time.Hour)},
	}
	s := suppress.New(rules)
	out := s.Apply(makeResults())
	found := false
	for _, r := range out {
		if r.Service == "api" {
			found = true
		}
	}
	if !found {
		t.Error("api should not be suppressed by an expired rule")
	}
}

func TestApply_CaseInsensitiveServiceName(t *testing.T) {
	rules := []suppress.Rule{
		{Service: "API"},
	}
	s := suppress.New(rules)
	out := s.Apply(makeResults())
	for _, r := range out {
		if r.Service == "api" {
			t.Error("api should be suppressed by case-insensitive rule")
		}
	}
}

func TestIsExpired_ZeroTime_NeverExpires(t *testing.T) {
	r := suppress.Rule{Service: "svc"}
	if r.IsExpired(time.Now()) {
		t.Error("rule with zero ExpiresAt should never expire")
	}
}

func TestIsExpired_FutureTime_NotExpired(t *testing.T) {
	r := suppress.Rule{Service: "svc", ExpiresAt: time.Now().Add(time.Hour)}
	if r.IsExpired(time.Now()) {
		t.Error("rule with future expiry should not be expired")
	}
}
