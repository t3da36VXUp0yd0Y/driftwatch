package suppress

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteRules writes a human-readable summary of the active suppression rules
// to w. Expired rules are omitted.
func WriteRules(w io.Writer, rules []Rule) error {
	now := time.Now()
	active := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if !r.IsExpired(now) {
			active = append(active, r)
		}
	}

	if len(active) == 0 {
		_, err := fmt.Fprintln(w, "no active suppression rules")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tREASON\tEXPIRES")
	for _, r := range active {
		reason := r.Reason
		if reason == "" {
			reason = "(all)"
		}
		expires := "never"
		if !r.ExpiresAt.IsZero() {
			expires = r.ExpiresAt.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.Service, reason, expires)
	}
	return tw.Flush()
}
