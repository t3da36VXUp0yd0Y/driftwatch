// Package cache implements a lightweight in-memory TTL cache used by the
// driftwatch runner to store recent drift detection results. Caching reduces
// redundant Docker API calls when the scheduler fires at a high frequency.
//
// Entries are keyed by an arbitrary string (typically a service name or a
// composite run key) and expire automatically after a configurable duration.
// A TTL of zero disables the cache entirely, which is the safe default for
// one-shot CLI invocations.
package cache
