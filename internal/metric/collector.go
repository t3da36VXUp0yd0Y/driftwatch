package metric

import (
	"github.com/user/driftwatch/internal/drift"
)

// Well-known counter names.
const (
	KeyChecksTotal  = "checks_total"
	KeyDriftTotal   = "drift_total"
	KeyHealthy      = "healthy_total"
	KeyErrors       = "errors_total"
)

// Collector wraps Counter and provides drift-aware recording helpers.
type Collector struct {
	*Counter
}

// NewCollector returns a Collector backed by a new Counter.
func NewCollector() *Collector {
	return &Collector{Counter: New()}
}

// RecordResults updates counters from a slice of drift results.
func (c *Collector) RecordResults(results []drift.Result) {
	for _, r := range results {
		c.Inc(KeyChecksTotal)
		if r.Drifted {
			c.Inc(KeyDriftTotal)
		} else {
			c.Inc(KeyHealthy)
		}
	}
}

// RecordError increments the error counter.
func (c *Collector) RecordError() {
	c.Inc(KeyErrors)
}
