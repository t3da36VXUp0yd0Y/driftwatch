// Package cache provides a simple in-memory TTL cache for storing
// drift detection results to avoid redundant Docker API calls.
package cache

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry holds a cached drift result along with its expiry time.
type Entry struct {
	Results   []drift.Result
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory store for drift results keyed by service name.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a new Cache with the given TTL duration.
// A zero TTL disables caching (entries are never considered valid).
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores results for the given key, expiring after the configured TTL.
func (c *Cache) Set(key string, results []drift.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = Entry{
		Results:   results,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves results for the given key. Returns the results and true if the
// entry exists and has not expired; otherwise returns nil and false.
func (c *Cache) Get(key string) ([]drift.Result, bool) {
	if c.ttl == 0 {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Results, true
}

// Invalidate removes the entry for the given key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]Entry)
}
