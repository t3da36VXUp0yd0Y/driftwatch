package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Report holds the results of a drift detection run.
type Report struct {
	GeneratedAt time.Time
	Results     []drift.Result
}

// New creates a new Report from drift detection results.
func New(results []drift.Result) *Report {
	return &Report{
		GeneratedAt: time.Now(),
		Results:     results,
	}
}

// HasDrift returns true if any result in the report indicates drift.
func (r *Report) HasDrift() bool {
	for _, res := range r.Results {
		if drift.HasDrift(res) {
			return true
		}
	}
	return false
}

// Write writes the report in the given format to the provided writer.
func (r *Report) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, r)
	default:
		return writeText(w, r)
	}
}

// Print writes the report to stdout.
func (r *Report) Print(format Format) error {
	return r.Write(os.Stdout, format)
}

func writeText(w io.Writer, r *Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Drift Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "%-20s\t%-8s\t%-40s\t%-40s\n", "SERVICE", "STATUS", "EXPECTED", "ACTUAL")
	fmt.Fprintf(tw, "%-20s\t%-8s\t%-40s\t%-40s\n", "-------", "------", "--------", "------")
	for _, res := range r.Results {
		status := "OK"
		if drift.HasDrift(res) {
			status = "DRIFT"
		}
		fmt.Fprintf(tw, "%-20s\t%-8s\t%-40s\t%-40s\n",
			res.ServiceName, status, res.ExpectedImage, res.ActualImage)
	}
	return tw.Flush()
}
