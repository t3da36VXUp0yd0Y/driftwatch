// Package throttle provides a simple time-based throttle to suppress
// repeated notifications or actions within a cooldown window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last trigger time per key and suppresses calls
// that occur within the cooldown duration.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New returns a Throttle with the given cooldown duration.
// Returns an error if cooldown is zero or negative.
func New(cooldown time.Duration) (*Throttle, error) {
	if cooldown <= 0 {
		return nil, ErrInvalidCooldown
	}
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}, nil
}

// ErrInvalidCooldown is returned when a non-positive cooldown is provided.
var ErrInvalidCooldown = throttleError("cooldown must be greater than zero")

type throttleError string

func (e throttleError) Error() string { return string(e) }

// Allow returns true if the key has not been triggered within the cooldown
// window, and records the current time as the last trigger.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the trigger history for the given key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// ResetAll clears all trigger history.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}
