// Package exit defines exit codes used by driftwatch to communicate
// the outcome of a drift check to the calling shell or CI system.
package exit

import "os"

// Code represents a process exit code.
type Code int

const (
	// OK indicates no drift was detected and all services matched their
	// declared configuration.
	OK Code = 0

	// DriftDetected indicates that at least one service was found to have
	// drifted from its declared configuration.
	DriftDetected Code = 1

	// ConfigError indicates the tool could not load or parse its
	// configuration file.
	ConfigError Code = 2

	// ClientError indicates a failure communicating with the Docker daemon
	// or another external dependency.
	ClientError Code = 3
)

// String returns a human-readable label for the exit code.
func (c Code) String() string {
	switch c {
	case OK:
		return "ok"
	case DriftDetected:
		return "drift_detected"
	case ConfigError:
		return "config_error"
	case ClientError:
		return "client_error"
	default:
		return "unknown"
	}
}

// Exit terminates the process with the given Code.
func Exit(c Code) {
	os.Exit(int(c))
}
