// Package dedupe provides deduplication of drift results within a sliding
// window, suppressing repeated alerts for the same drifted service.
package dedupe

import (
	"sync"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Entry records the last time a drift result was seen for a service.
type Entry struct {
	Service  string
	LastSeen time.Time
}

// Deduplicator suppresses duplicate drift results within a configurable TTL.
type Deduplicator struct {
	mu      sync.Mutex
	ttl     time.Duration
	seen    map[string]time.Time
	nowFunc func() time.Time
}

// New creates a Deduplicator with the given TTL. Returns an error if ttl <= 0.
func New(ttl time.Duration) (*Deduplicator, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Deduplicator{
		ttl:     ttl,
		seen:    make(map[string]time.Time),
		nowFunc: time.Now,
	}, nil
}

// ErrInvalidTTL is returned when a non-positive TTL is provided.
var ErrInvalidTTL = errInvalidTTL("dedupe: TTL must be greater than zero")

type errInvalidTTL string

func (e errInvalidTTL) Error() string { return string(e) }

// Filter returns only those results that have not been seen within the TTL.
// Results with no drift are always passed through unchanged.
func (d *Deduplicator) Filter(results []drift.Result) []drift.Result {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	out := make([]drift.Result, 0, len(results))

	for _, r := range results {
		if !r.Drifted {
			out = append(out, r)
			continue
		}
		last, exists := d.seen[r.Service]
		if !exists || now.Sub(last) >= d.ttl {
			d.seen[r.Service] = now
			out = append(out, r)
		}
	}
	return out
}

// Flush removes all recorded entries, resetting deduplication state.
func (d *Deduplicator) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

// Entries returns a snapshot of all currently tracked service entries.
func (d *Deduplicator) Entries() []Entry {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]Entry, 0, len(d.seen))
	for svc, t := range d.seen {
		out = append(out, Entry{Service: svc, LastSeen: t})
	}
	return out
}
