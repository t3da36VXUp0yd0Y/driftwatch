package report

import (
	"encoding/json"
	"io"

	"github.com/driftwatch/internal/drift"
)

type jsonReport struct {
	GeneratedAt string        `json:"generated_at"`
	HasDrift    bool          `json:"has_drift"`
	Results     []jsonResult  `json:"results"`
}

type jsonResult struct {
	ServiceName   string `json:"service_name"`
	ExpectedImage string `json:"expected_image"`
	ActualImage   string `json:"actual_image"`
	Running       bool   `json:"running"`
	Drifted       bool   `json:"drifted"`
}

func writeJSON(w io.Writer, r *Report) error {
	results := make([]jsonResult, len(r.Results))
	for i, res := range r.Results {
		results[i] = jsonResult{
			ServiceName:   res.ServiceName,
			ExpectedImage: res.ExpectedImage,
			ActualImage:   res.ActualImage,
			Running:       res.Running,
			Drifted:       drift.HasDrift(res),
		}
	}

	payload := jsonReport{
		GeneratedAt: r.GeneratedAt.UTC().Format("2006-01-02T15:04:05Z"),
		HasDrift:    r.HasDrift(),
		Results:     results,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
