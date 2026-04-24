package pipeline

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/drift"
)

// FuncStage wraps a plain function as a Stage, useful for inline or
// test-only stages without defining a full type.
type FuncStage struct {
	name string
	fn   func(ctx context.Context, results []drift.Result) ([]drift.Result, error)
}

// NewFuncStage creates a FuncStage with the given name and function.
// Returns an error if name is empty or fn is nil.
func NewFuncStage(name string, fn func(ctx context.Context, results []drift.Result) ([]drift.Result, error)) (*FuncStage, error) {
	if name == "" {
		return nil, fmt.Errorf("pipeline: stage name must not be empty")
	}
	if fn == nil {
		return nil, fmt.Errorf("pipeline: stage function must not be nil")
	}
	return &FuncStage{name: name, fn: fn}, nil
}

// Name returns the stage name.
func (f *FuncStage) Name() string { return f.name }

// Run executes the wrapped function.
func (f *FuncStage) Run(ctx context.Context, results []drift.Result) ([]drift.Result, error) {
	return f.fn(ctx, results)
}

// NoopStage is a pass-through stage that returns results unchanged.
// It is useful as a placeholder or in tests.
type NoopStage struct{ label string }

// NewNoopStage creates a NoopStage with the given label.
func NewNoopStage(label string) *NoopStage { return &NoopStage{label: label} }

// Name returns the noop stage label.
func (n *NoopStage) Name() string { return n.label }

// Run returns the input results unchanged.
func (n *NoopStage) Run(_ context.Context, results []drift.Result) ([]drift.Result, error) {
	return results, nil
}
