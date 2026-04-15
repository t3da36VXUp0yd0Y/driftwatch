package runner

import (
	"context"
	"fmt"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/report"
)

// Runner orchestrates the full drift detection pipeline.
type Runner struct {
	cfg    *config.Config
	client DockerClient
}

// DockerClient is the interface Runner depends on for container inspection.
type DockerClient interface {
	GetContainerInfo(ctx context.Context, name string) (*docker.ContainerInfo, error)
}

// New creates a Runner with the given config and docker client.
func New(cfg *config.Config, client DockerClient) *Runner {
	return &Runner{cfg: cfg, client: client}
}

// Run executes drift detection for all configured services and returns a Report.
func (r *Runner) Run(ctx context.Context, format report.Format) (*report.Report, error) {
	results := make([]drift.Result, 0, len(r.cfg.Services))

	for _, svc := range r.cfg.Services {
		info, err := r.client.GetContainerInfo(ctx, svc.Container)
		if err != nil {
			// Treat fetch error as service not running.
			results = append(results, drift.Result{
				Service: svc.Name,
				Drifted: true,
				Reason:  fmt.Sprintf("could not inspect container: %v", err),
			})
			continue
		}

		result := drift.Detect(svc, info)
		results = append(results, result)
	}

	rep, err := report.New(results, format)
	if err != nil {
		return nil, fmt.Errorf("building report: %w", err)
	}

	return rep, nil
}
