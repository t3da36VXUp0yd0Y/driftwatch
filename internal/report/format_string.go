package report

import "fmt"

// ParseFormat parses a string into a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown report format %q: must be one of [text, json]", s)
	}
}

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}
