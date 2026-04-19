// Package health provides container health status checking for drift detection.
package health

import (
	"fmt"
	"io"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Status represents the health state of a container.
type Status int

const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusUnhealthy
	StatusStarting
)

// String returns a human-readable health status.
func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	case StatusStarting:
		return "starting"
	default:
		return "unknown"
	}
}

// Result holds the health check outcome for a single service.
type Result struct {
	Service string
	Status  Status
	Detail  string
}

// Check evaluates the health status of each drift result.
// Services that are not running are marked unhealthy.
func Check(results []drift.Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		hr := Result{Service: r.Service}
		if !r.Running {
			hr.Status = StatusUnhealthy
			hr.Detail = "container is not running"
		} else if r.Drifted {
			hr.Status = StatusUnhealthy
			hr.Detail = fmt.Sprintf("drift detected: %s", r.Reason)
		} else {
			hr.Status = StatusHealthy
			hr.Detail = "ok"
		}
		out = append(out, hr)
	}
	return out
}

// Write prints health results to w in a human-readable format.
func Write(w io.Writer, results []Result) {
	for _, r := range results {
		fmt.Fprintf(w, "%-30s %s\n", r.Service, r.Status)
		if r.Detail != "ok" && r.Detail != "" {
			fmt.Fprintf(w, "  detail: %s\n", r.Detail)
		}
	}
}
