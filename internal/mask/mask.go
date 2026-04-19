// Package mask provides field-level masking for drift results before output.
package mask

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Rule defines a masking rule for a specific field pattern.
type Rule struct {
	Field       string
	Replacement string
}

// Masker applies masking rules to drift results.
type Masker struct {
	rules []Rule
}

// DefaultReplacement is used when a rule has no explicit replacement.
const DefaultReplacement = "***"

// New creates a Masker with the given rules.
// Rules with empty Field values are ignored.
func New(rules []Rule) *Masker {
	filtered := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			continue
		}
		if r.Replacement == "" {
			r.Replacement = DefaultReplacement
		}
		filtered = append(filtered, r)
	}
	return &Masker{rules: filtered}
}

// Apply returns a copy of results with sensitive fields masked.
func (m *Masker) Apply(results []drift.Result) []drift.Result {
	if len(m.rules) == 0 {
		return results
	}
	out := make([]drift.Result, len(results))
	for i, r := range results {
		out[i] = r
		out[i].Expected = m.maskValue(r.Field, r.Expected)
		out[i].Actual = m.maskValue(r.Field, r.Actual)
	}
	return out
}

func (m *Masker) maskValue(field, value string) string {
	for _, rule := range m.rules {
		if strings.EqualFold(field, rule.Field) {
			return rule.Replacement
		}
	}
	return value
}
