package pipeline_test

import (
	"context"
	"testing"

	"github.com/driftwatch/internal/pipeline"
)

func TestNewFuncStage_EmptyName(t *testing.T) {
	_, err := pipeline.NewFuncStage("", func(_ context.Context, r interface{}) (interface{}, error) { return r, nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewFuncStage_NilFunc(t *testing.T) {
	_, err := pipeline.NewFuncStage("valid", nil)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}

func TestFuncStage_Name(t *testing.T) {
	s, err := pipeline.NewFuncStage("my-stage", func(ctx context.Context, r interface{}) (interface{}, error) {
		return r, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name() != "my-stage" {
		t.Fatalf("expected name %q, got %q", "my-stage", s.Name())
	}
}

func TestNoopStage_Name(t *testing.T) {
	n := pipeline.NewNoopStage("passthrough")
	if n.Name() != "passthrough" {
		t.Fatalf("expected %q, got %q", "passthrough", n.Name())
	}
}

func TestNoopStage_ReturnsInputUnchanged(t *testing.T) {
	n := pipeline.NewNoopStage("noop")
	in := makeResults("alpha", "beta", "gamma")
	out, err := n.Run(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(in) {
		t.Fatalf("expected %d results, got %d", len(in), len(out))
	}
	for i, r := range out {
		if r.Service != in[i].Service {
			t.Fatalf("result %d: expected service %q, got %q", i, in[i].Service, r.Service)
		}
	}
}
