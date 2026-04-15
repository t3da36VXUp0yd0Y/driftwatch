package runner_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/report"
	"github.com/user/driftwatch/internal/runner"
)

type mockClient struct {
	containers map[string]*docker.ContainerInfo
	err        error
}

func (m *mockClient) GetContainerInfo(_ context.Context, name string) (*docker.ContainerInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	info, ok := m.containers[name]
	if !ok {
		return nil, errors.New("container not found")
	}
	return info, nil
}

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestRun_NoDrift(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "web", Container: "web", Image: "nginx:latest"},
	})
	client := &mockClient{
		containers: map[string]*docker.ContainerInfo{
			"web": {Name: "web", Image: "nginx:latest", Running: true},
		},
	}
	r := runner.New(cfg, client)
	rep, err := r.Run(context.Background(), report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestRun_ClientError(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Container: "api", Image: "myapp:v1"},
	})
	client := &mockClient{err: errors.New("daemon unreachable")}
	r := runner.New(cfg, client)
	rep, err := r.Run(context.Background(), report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rep.HasDrift() {
		t.Error("expected drift when client errors")
	}
}

func TestRun_ImageMismatch(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "db", Container: "db", Image: "postgres:15"},
	})
	client := &mockClient{
		containers: map[string]*docker.ContainerInfo{
			"db": {Name: "db", Image: "postgres:14", Running: true},
		},
	}
	r := runner.New(cfg, client)
	rep, err := r.Run(context.Background(), report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rep.HasDrift() {
		t.Error("expected drift on image mismatch")
	}
}
