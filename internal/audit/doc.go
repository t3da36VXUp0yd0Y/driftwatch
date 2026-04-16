// Package audit provides an append-only audit log for drift check events.
// Each call to Record appends one JSON entry per service result to the
// configured log file, enabling historical review of when and where
// configuration drift was detected.
package audit
