// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently drift checks are allowed to run.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when a call is rejected due to rate limiting.
var ErrRateLimited = errors.New("rate limited: too many requests")

// Limiter controls the rate of operations using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   int
	max      int
	refillAt time.Time
	window   time.Duration
}

// New creates a Limiter that allows max calls per window duration.
// Returns an error if max < 1 or window <= 0.
func New(max int, window time.Duration) (*Limiter, error) {
	if max < 1 {
		return nil, errors.New("max must be at least 1")
	}
	if window <= 0 {
		return nil, errors.New("window must be positive")
	}
	return &Limiter{
		tokens:   max,
		max:      max,
		window:   window,
		refillAt: time.Now().Add(window),
	}, nil
}

// Allow returns nil if the call is permitted, or ErrRateLimited if the
// bucket is empty. Tokens are refilled after each window elapses.
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.After(l.refillAt) {
		l.tokens = l.max
		l.refillAt = now.Add(l.window)
	}

	if l.tokens <= 0 {
		return ErrRateLimited
	}
	l.tokens--
	return nil
}

// Remaining returns the number of tokens left in the current window.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
