// Package notify implements lightweight notification hooks for driftwatch.
//
// A Notifier writes human-readable drift alerts to any io.Writer (stdout,
// a file, or a network connection).  The notification level controls whether
// all results are emitted or only those that contain drift.
//
// Basic usage:
//
//	n := notify.New(os.Stderr, notify.LevelDriftOnly)
//	count, err := n.Notify(results)
package notify
