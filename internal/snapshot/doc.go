// Package snapshot captures and persists point-in-time views of drift
// detection results. Snapshots can be saved to disk as JSON files and
// later loaded for comparison or audit purposes.
//
// Usage:
//
//	s := snapshot.New(results)
//	s.Meta["environment"] = "production"
//	path, err := snapshot.Save("/var/lib/driftwatch/snapshots", s)
package snapshot
