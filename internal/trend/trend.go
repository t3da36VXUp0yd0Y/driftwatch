// Package trend analyses drift results over time to identify recurring or worsening services.
package trend

import (
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a single historical observation.
type Entry struct {
	Time    time.Time
	Results []drift.Result
}

// ServiceTrend summarises drift frequency for a single service.
type ServiceTrend struct {
	Service      string
	ObservedAt   []time.Time
	DriftCount   int
	HealthyCount int
}

// Analyse examines a slice of historical entries and returns per-service trends.
func Analyse(entries []Entry) []ServiceTrend {
	index := map[string]*ServiceTrend{}

	for _, e := range entries {
		for _, r := range e.Results {
			st, ok := index[r.Service]
			if !ok {
				st = &ServiceTrend{Service: r.Service}
				index[r.Service] = st
			}
			st.ObservedAt = append(st.ObservedAt, e.Time)
			if r.Drifted {
				st.DriftCount++
			} else {
				st.HealthyCount++
			}
		}
	}

	out := make([]ServiceTrend, 0, len(index))
	for _, st := range index {
		out = append(out, *st)
	}
	return out
}

// DriftRate returns the fraction of observations where the service was drifted.
// Returns 0 if there are no observations.
func DriftRate(st ServiceTrend) float64 {
	total := st.DriftCount + st.HealthyCount
	if total == 0 {
		return 0
	}
	return float64(st.DriftCount) / float64(total)
}
