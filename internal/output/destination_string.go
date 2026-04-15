package output

import "fmt"

// ParseDestination converts a string flag value into a Destination constant.
// Accepted values: "stdout", "file".
func ParseDestination(s string) (Destination, error) {
	switch s {
	case "stdout":
		return Stdout, nil
	case "file":
		return File, nil
	default:
		return Stdout, fmt.Errorf("output: unknown destination %q, must be one of: stdout, file", s)
	}
}

// String returns the canonical string representation of a Destination.
func (d Destination) String() string {
	switch d {
	case Stdout:
		return "stdout"
	case File:
		return "file"
	default:
		return fmt.Sprintf("Destination(%d)", int(d))
	}
}
