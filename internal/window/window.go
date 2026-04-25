// Package window provides a sliding time-window aggregator for drift results.
// It accumulates check results over a configurable duration and exposes
// summary statistics useful for trend analysis and alerting.
package window

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry holds a single recorded observation.
type Entry struct {
	Timestamp time.Time
	Results   []drift.Result
}

// Window is a sliding time-window that retains drift results.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	entries  []Entry
}

// New creates a Window with the given retention duration.
// Returns an error if duration is zero or negative.
func New(d time.Duration) (*Window, error) {
	if d <= 0 {
		return nil, ErrInvalidDuration
	}
	return &Window{duration: d}, nil
}

// Add appends a new observation and evicts entries outside the window.
func (w *Window) Add(results []drift.Result) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = append(w.entries, Entry{
		Timestamp: time.Now(),
		Results:   results,
	})
	w.evict()
}

// Entries returns a copy of all observations currently in the window.
func (w *Window) Entries() []Entry {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	out := make([]Entry, len(w.entries))
	copy(out, w.entries)
	return out
}

// DriftedCount returns the number of drifted results across all observations
// currently inside the window.
func (w *Window) DriftedCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	count := 0
	for _, e := range w.entries {
		for _, r := range e.Results {
			if r.Drifted {
				count++
			}
		}
	}
	return count
}

// evict removes entries older than the window duration. Must be called with mu held.
func (w *Window) evict() {
	cutoff := time.Now().Add(-w.duration)
	i := 0
	for i < len(w.entries) && w.entries[i].Timestamp.Before(cutoff) {
		i++
	}
	w.entries = w.entries[i:]
}
