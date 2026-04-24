package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/pipeline"
)

func makeResults(names ...string) []drift.Result {
	out := make([]drift.Result, len(names))
	for i, n := range names {
		out[i] = drift.Result{Service: n}
	}
	return out
}

func TestNew_NilStageReturnsError(t *testing.T) {
	_, err := pipeline.New(nil)
	if err == nil {
		t.Fatal("expected error for nil stage, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	p, err := pipeline.New(pipeline.NewNoopStage("a"), pipeline.NewNoopStage("b"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := p.Stages()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected stage names: %v", names)
	}
}

func TestRun_PassesThroughNoop(t *testing.T) {
	p, _ := pipeline.New(pipeline.NewNoopStage("noop"))
	in := makeResults("svc-a", "svc-b")
	out, err := p.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(in) {
		t.Fatalf("expected %d results, got %d", len(in), len(out))
	}
}

func TestRun_StageError_Propagates(t *testing.T) {
	errBoom := errors.New("boom")
	fail, _ := pipeline.NewFuncStage("fail", func(_ context.Context, r []drift.Result) ([]drift.Result, error) {
		return nil, errBoom
	})
	p, _ := pipeline.New(fail)
	_, err := p.Run(context.Background(), makeResults("svc-a"))
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
}

func TestRun_StagesExecuteInOrder(t *testing.T) {
	var order []string
	mkStage := func(name string) pipeline.Stage {
		s, _ := pipeline.NewFuncStage(name, func(_ context.Context, r []drift.Result) ([]drift.Result, error) {
			order = append(order, name)
			return r, nil
		})
		return s
	}
	p, _ := pipeline.New(mkStage("first"), mkStage("second"), mkStage("third"))
	p.Run(context.Background(), makeResults("svc"))
	if len(order) != 3 || order[0] != "first" || order[1] != "second" || order[2] != "third" {
		t.Fatalf("unexpected execution order: %v", order)
	}
}

func TestRun_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	noop := pipeline.NewNoopStage("noop")
	p, _ := pipeline.New(noop)
	_, err := p.Run(ctx, makeResults("svc"))
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
