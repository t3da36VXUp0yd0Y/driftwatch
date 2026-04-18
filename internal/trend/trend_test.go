package trend_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/trend"
)

func makeEntries() []trend.Entry {
	t0 := time.Now().Add(-2 * time.Hour)
	t1 := time.Now().Add(-1 * time.Hour)
	return []trend.Entry{
		{
			Time: t0,
			Results: []drift.Result{
				{Service: "api", Drifted: true},
				{Service: "worker", Drifted: false},
			},
		},
		{
			Time: t1,
			Results: []drift.Result{
				{Service: "api", Drifted: true},
				{Service: "worker", Drifted: true},
			},
		},
	}
}

func TestAnalyse_ReturnsTrendsPerService(t *testing.T) {
	entries := makeEntries()
	trends := trend.Analyse(entries)
	if len(trends) != 2 {
		t.Fatalf("expected 2 trends, got %d", len(trends))
	}
}

func TestAnalyse_DriftCounts(t *testing.T) {
	entries := makeEntries()
	trends := trend.Analyse(entries)

	idx := map[string]trend.ServiceTrend{}
	for _, st := range trends {
		idx[st.Service] = st
	}

	if idx["api"].DriftCount != 2 {
		t.Errorf("api: expected DriftCount 2, got %d", idx["api"].DriftCount)
	}
	if idx["worker"].DriftCount != 1 {
		t.Errorf("worker: expected DriftCount 1, got %d", idx["worker"].DriftCount)
	}
	if idx["worker"].HealthyCount != 1 {
		t.Errorf("worker: expected HealthyCount 1, got %d", idx["worker"].HealthyCount)
	}
}

func TestDriftRate_Full(t *testing.T) {
	st := trend.ServiceTrend{Service: "api", DriftCount: 3, HealthyCount: 0}
	if rate := trend.DriftRate(st); rate != 1.0 {
		t.Errorf("expected 1.0, got %f", rate)
	}
}

func TestDriftRate_Zero(t *testing.T) {
	st := trend.ServiceTrend{Service: "api", DriftCount: 0, HealthyCount: 0}
	if rate := trend.DriftRate(st); rate != 0 {
		t.Errorf("expected 0, got %f", rate)
	}
}

func TestAnalyse_EmptyEntries(t *testing.T) {
	trends := trend.Analyse([]trend.Entry{})
	if len(trends) != 0 {
		t.Errorf("expected empty trends, got %d", len(trends))
	}
}
