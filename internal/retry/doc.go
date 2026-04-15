// Package retry provides a lightweight retry helper used by driftwatch
// when communicating with the Docker daemon or other external services.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return client.Ping(ctx)
//	})
package retry
