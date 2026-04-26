// Package compare provides field-level comparison between expected and actual
// service configurations, producing a structured diff suitable for reporting.
package compare

import (
	"fmt"
	"strings"
)

// Field represents a single compared configuration field.
type Field struct {
	Name     string
	Expected string
	Actual   string
}

// Changed reports whether the expected and actual values differ.
func (f Field) Changed() bool {
	return f.Expected != f.Actual
}

// String returns a human-readable description of the field comparison.
func (f Field) String() string {
	if f.Changed() {
		return fmt.Sprintf("%s: expected %q, got %q", f.Name, f.Expected, f.Actual)
	}
	return fmt.Sprintf("%s: %q (unchanged)", f.Name, f.Actual)
}

// Result holds the comparison outcome for a single service.
type Result struct {
	Service string
	Fields  []Field
}

// HasDrift reports whether any field in the result has changed.
func (r Result) HasDrift() bool {
	for _, f := range r.Fields {
		if f.Changed() {
			return true
		}
	}
	return false
}

// DriftedFields returns only the fields that have changed.
func (r Result) DriftedFields() []Field {
	out := make([]Field, 0, len(r.Fields))
	for _, f := range r.Fields {
		if f.Changed() {
			out = append(out, f)
		}
	}
	return out
}

// Summary returns a short human-readable summary of drifted fields.
func (r Result) Summary() string {
	drifted := r.DriftedFields()
	if len(drifted) == 0 {
		return fmt.Sprintf("%s: no drift detected", r.Service)
	}
	parts := make([]string, len(drifted))
	for i, f := range drifted {
		parts[i] = f.String()
	}
	return fmt.Sprintf("%s: %s", r.Service, strings.Join(parts, "; "))
}

// Pair holds an expected/actual string pair for a named field.
type Pair struct {
	Name     string
	Expected string
	Actual   string
}

// Compare builds a Result for the given service by evaluating each Pair.
// Pairs with empty names are silently skipped.
func Compare(service string, pairs []Pair) Result {
	fields := make([]Field, 0, len(pairs))
	for _, p := range pairs {
		if p.Name == "" {
			continue
		}
		fields = append(fields, Field{
			Name:     p.Name,
			Expected: p.Expected,
			Actual:   p.Actual,
		})
	}
	return Result{
		Service: service,
		Fields:  fields,
	}
}
