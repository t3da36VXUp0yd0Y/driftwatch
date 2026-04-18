// Package digest computes and compares fingerprints for drift results,
// allowing callers to detect whether a result set has changed since last check.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Compute returns a stable SHA-256 hex fingerprint for the given results.
// Results are sorted by service name before hashing so order is irrelevant.
func Compute(results []drift.Result) string {
	type entry struct {
		service string
		drifted bool
		reason  string
	}

	entries := make([]entry, len(results))
	for i, r := range results {
		entries[i] = entry{r.Service, r.Drifted, r.Reason}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].service < entries[j].service
	})

	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s:%v:%s|", e.service, e.drifted, e.reason)
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// Changed returns true when the fingerprint of current differs from previous.
func Changed(previous string, current []drift.Result) bool {
	return previous != Compute(current)
}
