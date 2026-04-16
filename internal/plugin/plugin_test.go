package plugin_test

import (
	"errors"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/driftwatch/driftwatch/internal/plugin"
)

func makeResult(service string, drifted bool) drift.Result {
	return drift.Result{Service: service, Drifted: drifted}
}

func TestRegister_And_Names(t *testing.T) {
	r := plugin.New()
	_ = r.Register("b", func() ([]drift.Result, error) { return nil, nil })
	_ = r.Register("a", func() ([]drift.Result, error) { return nil, nil })
	names := r.Names()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	r := plugin.New()
	_ = r.Register("x", func() ([]drift.Result, error) { return nil, nil })
	err := r.Register("x", func() ([]drift.Result, error) { return nil, nil })
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRun_NotFound(t *testing.T) {
	r := plugin.New()
	_, err := r.Run("missing")
	if err == nil {
		t.Fatal("expected error for missing plugin")
	}
}

func TestRun_ReturnsResults(t *testing.T) {
	r := plugin.New()
	_ = r.Register("demo", func() ([]drift.Result, error) {
		return []drift.Result{makeResult("svc", true)}, nil
	})
	res, err := r.Run("demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Service != "svc" {
		t.Errorf("unexpected results: %v", res)
	}
}

func TestRunAll_MergesResults(t *testing.T) {
	r := plugin.New()
	_ = r.Register("p1", func() ([]drift.Result, error) {
		return []drift.Result{makeResult("a", false)}, nil
	})
	_ = r.Register("p2", func() ([]drift.Result, error) {
		return []drift.Result{makeResult("b", true)}, nil
	})
	res, err := r.RunAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 results, got %d", len(res))
	}
}

func TestRunAll_PropagatesError(t *testing.T) {
	r := plugin.New()
	_ = r.Register("bad", func() ([]drift.Result, error) {
		return nil, errors.New("boom")
	})
	_, err := r.RunAll()
	if err == nil {
		t.Fatal("expected error from failing plugin")
	}
}
