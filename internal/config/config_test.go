package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/you/driftwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `
version: "1"
services:
  - name: api
    image: myrepo/api:latest
    replicas: 2
    environment:
      LOG_LEVEL: info
    ports:
      - "8080:8080"
`
	p := writeTemp(t, raw)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	svc := cfg.Services[0]
	if svc.Name != "api" {
		t.Errorf("expected name 'api', got %q", svc.Name)
	}
	if svc.Replicas != 2 {
		t.Errorf("expected 2 replicas, got %d", svc.Replicas)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyServices(t *testing.T) {
	raw := `version: "1"
services: []
`
	p := writeTemp(t, raw)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for empty services")
	}
}

func TestLoad_MissingImage(t *testing.T) {
	raw := `version: "1"
services:
  - name: worker
    replicas: 1
`
	p := writeTemp(t, raw)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing image")
	}
}
