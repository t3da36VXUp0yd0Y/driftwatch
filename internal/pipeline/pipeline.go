// Package pipeline provides a composable stage-based processing pipeline
// for drift detection results. Stages are applied in order, and each stage
// may transform, filter, or annotate the result set.
package pipeline

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// Stage is a single processing step in the pipeline.
type Stage interface {
	// Name returns a human-readable identifier for the stage.
	Name() string
	// Run processes the results and returns a (possibly modified) slice.
	Run(ctx context.Context, results []drift.Result) ([]drift.Result, error)
}

// Pipeline executes an ordered sequence of stages against drift results.
type Pipeline struct {
	stages []Stage
}

// New creates a Pipeline with the provided stages.
func New(stages ...Stage) (*Pipeline, error) {
	for i, s := range stages {
		if s == nil {
			return nil, fmt.Errorf("pipeline: stage at index %d is nil", i)
		}
	}
	return &Pipeline{stages: stages}, nil
}

// Stages returns the names of all registered stages in order.
func (p *Pipeline) Stages() []string {
	names := make([]string, len(p.stages))
	for i, s := range p.stages {
		names[i] = s.Name()
	}
	return names
}

// Run executes each stage sequentially, passing the output of one stage
// as the input to the next. It returns the final result set or the first
// error encountered.
func (p *Pipeline) Run(ctx context.Context, results []drift.Result) ([]drift.Result, error) {
	current := results
	for _, s := range p.stages {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("pipeline: context cancelled before stage %q: %w", s.Name(), err)
		}
		var err error
		current, err = s.Run(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("pipeline: stage %q failed: %w", s.Name(), err)
		}
	}
	return current, nil
}
