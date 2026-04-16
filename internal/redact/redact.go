// Package redact provides utilities for sanitising sensitive values
// (e.g. image tags containing credentials) before they are written to
// any output or report.
package redact

import (
	"regexp"
	"strings"
)

// pattern matches an HTTP/HTTPS basic-auth prefix in a registry URL.
// e.g. https://user:password@registry.example.com/image:tag
var basicAuth = regexp.MustCompile(`(?i)(https?://)([^:]+:[^@]+@)`)

// Redactor holds compiled rules used to sanitise strings.
type Redactor struct {
	extra []*regexp.Regexp
}

// New returns a Redactor. Additional literal tokens to redact can be
// supplied; each will be replaced with "***".
func New(secrets ...string) *Redactor {
	r := &Redactor{}
	for _, s := range secrets {
		if s == "" {
			continue
		}
		r.extra = append(r.extra, regexp.MustCompile(regexp.QuoteMeta(s)))
	}
	return r
}

// Sanitise returns a copy of s with sensitive fragments replaced.
func (r *Redactor) Sanitise(s string) string {
	s = basicAuth.ReplaceAllString(s, "$1***@")
	for _, re := range r.extra {
		s = re.ReplaceAllString(s, "***")
	}
	return s
}

// SanitiseSlice applies Sanitise to every element of ss.
func (r *Redactor) SanitiseSlice(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = r.Sanitise(s)
	}
	return out
}

// ContainsSecret reports whether s contains any of the registered secrets.
func (r *Redactor) ContainsSecret(s string) bool {
	for _, re := range r.extra {
		if re.MatchString(s) {
			return true
		}
	}
	return strings.Contains(s, "@") && basicAuth.MatchString(s)
}
