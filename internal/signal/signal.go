// Package signal provides drift result aggregation by computing
// a weighted signal score across multiple check dimensions.
package signal

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Weight maps a check dimension name to its relative importance.
type Weight struct {
	Name   string
	Factor float64
}

// Score represents the computed signal for a single service.
type Score struct {
	Service string
	Value   float64
	Drifted bool
}

// DefaultWeights returns a sensible default weighting configuration.
func DefaultWeights() []Weight {
	return []Weight{
		{Name: "image", Factor: 1.0},
		{Name: "env", Factor: 0.8},
		{Name: "port", Factor: 0.6},
		{Name: "volume", Factor: 0.5},
		{Name: "health", Factor: 0.4},
	}
}

// Compute derives a signal score per service from drift results.
// Each drifted result contributes its weight factor to the total score.
// Results without a matching weight receive a default factor of 0.5.
func Compute(results []drift.Result, weights []Weight) []Score {
	wmap := buildWeightMap(weights)
	accum := make(map[string]float64)
	seen := make(map[string]bool)

	for _, r := range results {
		seen[r.Service] = true
		if !r.Drifted {
			continue
		}
		f, ok := wmap[r.Service]
		if !ok {
			f = 0.5
		}
		accum[r.Service] += f
	}

	out := make([]Score, 0, len(seen))
	for svc := range seen {
		out = append(out, Score{
			Service: svc,
			Value:   accum[svc],
			Drifted: accum[svc] > 0,
		})
	}
	return out
}

// Write renders signal scores to w in a human-readable format.
func Write(w io.Writer, scores []Score) {
	for _, s := range scores {
		status := "ok"
		if s.Drifted {
			status = "drifted"
		}
		fmt.Fprintf(w, "%-30s score=%.2f status=%s\n", s.Service, s.Value, status)
	}
}

func buildWeightMap(weights []Weight) map[string]float64 {
	m := make(map[string]float64, len(weights))
	for _, w := range weights {
		m[w.Name] = w.Factor
	}
	return m
}
