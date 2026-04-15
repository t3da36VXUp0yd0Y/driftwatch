// Package version exposes build-time metadata for the driftwatch CLI.
//
// Version, Commit, and Date are injected at compile time via -ldflags:
//
//	-ldflags "-X github.com/driftwatch/driftwatch/internal/version.Version=1.0.0"
//	         "-X github.com/driftwatch/driftwatch/internal/version.Commit=$(git rev-parse --short HEAD)"
//	         "-X github.com/driftwatch/driftwatch/internal/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
//
// When not set (e.g. during local development), sensible defaults are used.
package version
