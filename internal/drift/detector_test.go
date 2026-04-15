package drift

import (
	"testing"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/docker"
)

func TestDetect_NoDrift(t *testing.T) {
	services := []config.Service{
		{Name: "api", Image: "nginx:1.25"},
	}
	running := map[string]docker.ContainerInfo{
		"api": {Name: "api", Image: "nginx:1.25", State: "running"},
	}

	results := Detect(services, running)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got reason: %s", results[0].Reason)
	}
}

func TestDetect_ImageMismatch(t *testing.T) {
	services := []config.Service{
		{Name: "api", Image: "nginx:1.25"},
	}
	running := map[string]docker.ContainerInfo{
		"api": {Name: "api", Image: "nginx:1.24", State: "running"},
	}

	results := Detect(services, running)
	if !results[0].Drifted {
		t.Error("expected drift due to image mismatch")
	}
	if results[0].ActualImage != "nginx:1.24" {
		t.Errorf("unexpected ActualImage: %s", results[0].ActualImage)
	}
}

func TestDetect_ServiceNotRunning(t *testing.T) {
	services := []config.Service{
		{Name: "worker", Image: "myapp:latest"},
	}
	running := map[string]docker.ContainerInfo{}

	results := Detect(services, running)
	if !results[0].Drifted {
		t.Error("expected drift for missing service")
	}
	if results[0].Running {
		t.Error("expected Running=false")
	}
}

func TestHasDrift_True(t *testing.T) {
	results := []Result{{Drifted: false}, {Drifted: true}}
	if !HasDrift(results) {
		t.Error("expected HasDrift=true")
	}
}

func TestHasDrift_False(t *testing.T) {
	results := []Result{{Drifted: false}, {Drifted: false}}
	if HasDrift(results) {
		t.Error("expected HasDrift=false")
	}
}
