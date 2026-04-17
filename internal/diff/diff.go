// Package diff provides field-level diffing between expected and actual
// container state, producing human-readable change descriptions.
package diff

import "fmt"

// Field represents a single diffed field between expected and actual state.
type Field struct {
	Name     string
	Expected string
	Actual   string
}

// String returns a formatted description of the field difference.
func (f Field) String() string {
	return fmt.Sprintf("%s: expected %q, got %q", f.Name, f.Expected, f.Actual)
}

// Result holds all field-level differences for a single service.
type Result struct {
	Service string
	Fields  []Field
}

// HasDiff reports whether any field differences were found.
func (r Result) HasDiff() bool {
	return len(r.Fields) > 0
}

// Summary returns a short human-readable summary of the diff.
func (r Result) Summary() string {
	if !r.HasDiff() {
		return fmt.Sprintf("%s: no differences", r.Service)
	}
	return fmt.Sprintf("%s: %d field(s) differ", r.Service, len(r.Fields))
}

// Compare compares expected vs actual values for a set of named fields and
// returns a Result describing any differences found.
func Compare(service string, pairs map[string][2]string) Result {
	r := Result{Service: service}
	for name, pair := range pairs {
		if pair[0] != pair[1] {
			r.Fields = append(r.Fields, Field{
				Name:     name,
				Expected: pair[0],
				Actual:   pair[1],
			})
		}
	}
	return r
}
