package drift

import (
	"fmt"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/docker"
)

// Result describes the drift status of a single service.
type Result struct {
	ServiceName    string
	ExpectedImage  string
	ActualImage    string
	Running        bool
	Drifted        bool
	Reason         string
}

// Detect compares declared config against the actual running containers and
// returns a slice of Results, one per declared service.
func Detect(services []config.Service, running map[string]docker.ContainerInfo) []Result {
	results := make([]Result, 0, len(services))

	for _, svc := range services {
		r := Result{
			ServiceName:   svc.Name,
			ExpectedImage: svc.Image,
		}

		info, found := running[svc.Name]
		if !found {
			r.Running = false
			r.Drifted = true
			r.Reason = fmt.Sprintf("service %q is not running", svc.Name)
			results = append(results, r)
			continue
		}

		r.Running = true
		r.ActualImage = info.Image

		if info.Image != svc.Image {
			r.Drifted = true
			r.Reason = fmt.Sprintf("image mismatch: expected %q, got %q", svc.Image, info.Image)
		} else {
			r.Drifted = false
			r.Reason = "ok"
		}

		results = append(results, r)
	}

	return results
}

// HasDrift returns true if any result in the slice is drifted.
func HasDrift(results []Result) bool {
	for _, r := range results {
		if r.Drifted {
			return true
		}
	}
	return false
}
