// Package ratelimit provides a token-bucket rate limiter used to guard
// drift-check operations against excessive execution frequency.
//
// Usage:
//
//	limiter, err := ratelimit.New(10, time.Minute)
//	if err != nil { ... }
//	if err := limiter.Allow(); err != nil {
//		// back off or skip this cycle
//	}
package ratelimit
