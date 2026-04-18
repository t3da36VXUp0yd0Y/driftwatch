package metric

import (
	"sync"
	"time"
)

// Counter tracks named counters for drift check runs.
type Counter struct {
	mu      sync.Mutex
	counts  map[string]int64
	sampled time.Time
}

// New returns a new Counter.
func New() *Counter {
	return &Counter{
		counts:  make(map[string]int64),
		sampled: time.Now(),
	}
}

// Inc increments the named counter by 1.
func (c *Counter) Inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[name]++
}

// Add increments the named counter by n.
func (c *Counter) Add(name string, n int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[name] += n
}

// Get returns the current value of the named counter.
func (c *Counter) Get(name string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[name]
}

// Snapshot returns a copy of all counters and the time they were sampled.
func (c *Counter) Snapshot() (map[string]int64, time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]int64, len(c.counts))
	for k, v := range c.counts {
		copy[k] = v
	}
	c.sampled = time.Now()
	return copy, c.sampled
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]int64)
}
