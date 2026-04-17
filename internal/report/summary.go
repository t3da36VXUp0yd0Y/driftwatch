package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftwatch/internal/drift"
)

// Summary holds aggregated statistics about a drift detection run.
type Summary struct {
	Total   int
	Drifted int
	Healthy int
	Missing int
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

// DriftPercent returns the percentage of services that have drifted,
// rounded to two decimal places. Returns 0 if Total is zero.
func (s Summary) DriftPercent() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Drifted) / float64(s.Total) * 100
}

// WriteSummary writes a human-readable summary table to w.
func WriteSummary(w io.Writer, s Summary) error {
	line := strings.Repeat("-", 36)
	_, err := fmt.Fprintf(w,
		"%s\n  Drift Detection Summary\n%s\n  Total services : %d\n  Healthy        : %d\n  Drifted        : %d\n  Missing        : %d\n  Drift %%        : %.2f%%\n%s\n",
		line, line, s.Total, s.Healthy, s.Drifted, s.Missing, s.DriftPercent(), line,
	)
	return err
}
