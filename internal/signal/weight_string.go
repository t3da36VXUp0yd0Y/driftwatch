package signal

import (
	"fmt"
	"strings"
)

// String returns a human-readable representation of a Weight.
func (w Weight) String() string {
	return fmt.Sprintf("%s(%.2f)", w.Name, w.Factor)
}

// ParseWeight parses a "name=factor" string into a Weight.
// It returns an error if the format is invalid or the factor is not a
// valid positive number.
func ParseWeight(s string) (Weight, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Weight{}, fmt.Errorf("signal: invalid weight %q: expected name=factor", s)
	}
	var f float64
	if _, err := fmt.Sscanf(parts[1], "%f", &f); err != nil {
		return Weight{}, fmt.Errorf("signal: invalid factor in %q: %w", s, err)
	}
	if f < 0 {
		return Weight{}, fmt.Errorf("signal: factor must be non-negative, got %.2f", f)
	}
	return Weight{Name: parts[0], Factor: f}, nil
}
