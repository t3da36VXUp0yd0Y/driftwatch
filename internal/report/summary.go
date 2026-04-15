package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

// Summary holds aggregated statistics about a drift detection run.
type Summary struct {
	Total    int
	Drifted  int
	Healthy  int
	Missing  int
}

// NewSummary builds a Summary from a slice of DetectResult.
func NewSummary(results []drift.DetectResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch {
		case !r.Running:
			s.Missing++
			s.Drifted++
		case r.ImageDrift:
			s.Drifted++
		default:
			s.Healthy++
		}
	}
	return s
}

// WriteSummary writes a human-readable summary table to w.
func WriteSummary(w io.Writer, s Summary) error {
	line := strings.Repeat("-", 36)
	_, err := fmt.Fprintf(w,
		"%s\n  Drift Detection Summary\n%s\n  Total services : %d\n  Healthy        : %d\n  Drifted        : %d\n  Missing        : %d\n%s\n",
		line, line, s.Total, s.Healthy, s.Drifted, s.Missing, line,
	)
	return err
}
