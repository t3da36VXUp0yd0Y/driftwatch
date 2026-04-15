// Package alert provides threshold-based alerting for configuration drift results.
//
// It evaluates the number of drifted services against configurable warn and
// critical thresholds, returning a structured Alert that callers can act on
// or write to any io.Writer for human-readable output.
package alert
