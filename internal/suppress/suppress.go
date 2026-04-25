// Package suppress provides a mechanism to silence drift results for known,
// accepted deviations that should not trigger alerts or reports.
package suppress

import (
	"strings"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// Rule defines a suppression rule that matches a service and optional drift reason.
type Rule struct {
	// Service is the exact service name to suppress (case-insensitive).
	Service string
	// Reason is an optional substring match against the drift detail.
	// If empty, all drift for the service is suppressed.
	Reason string
	// ExpiresAt is the time after which the rule is no longer applied.
	// A zero value means the rule never expires.
	ExpiresAt time.Time
}

// IsExpired reports whether the rule has passed its expiry time.
func (r Rule) IsExpired(now time.Time) bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return now.After(r.ExpiresAt)
}

// Suppressor filters drift results according to a set of rules.
type Suppressor struct {
	rules []Rule
	now   func() time.Time
}

// New returns a Suppressor with the provided rules.
func New(rules []Rule) *Suppressor {
	return &Suppressor{
		rules: rules,
		now:   time.Now,
	}
}

// Apply returns only the results that are not matched by any active suppression rule.
func (s *Suppressor) Apply(results []drift.Result) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if !s.isSuppressed(r) {
			out = append(out, r)
		}
	}
	return out
}

// isSuppressed returns true if the result matches any active rule.
func (s *Suppressor) isSuppressed(r drift.Result) bool {
	now := s.now()
	for _, rule := range s.rules {
		if rule.IsExpired(now) {
			continue
		}
		if !strings.EqualFold(rule.Service, r.Service) {
			continue
		}
		if rule.Reason == "" {
			return true
		}
		if strings.Contains(r.Detail, rule.Reason) {
			return true
		}
	}
	return false
}
