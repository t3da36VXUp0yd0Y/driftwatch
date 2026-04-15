// Package runner wires together config loading, Docker inspection,
// drift detection, optional filtering, and report generation.
package runner

import (
	"context"
	"io"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/docker"
	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/filter"
	"github.com/example/driftwatch/internal/report"
)

// Runner orchestrates a single drift-check pass.
type Runner struct {
	cfg        *config.Config
	client     docker.Client
	filterOpts filter.Options
	reportOpts report.Options
}

// New constructs a Runner from the provided config path and options.
func New(cfgPath string, fo filter.Options, ro report.Options) (*Runner, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	return &Runner{cfg: cfg, client: cli, filterOpts: fo, reportOpts: ro}, nil
}

// Run executes the full drift-check pipeline and writes the report to w.
// It returns true when drift is detected in any service.
func (r *Runner) Run(ctx context.Context, w io.Writer) (bool, error) {
	running, err := r.client.ListContainers(ctx)
	if err != nil {
		return false, err
	}

	results := drift.Detect(r.cfg.Services, running)
	results = filter.Apply(results, r.filterOpts)

	rep, err := report.New(results, r.reportOpts)
	if err != nil {
		return false, err
	}

	if err := rep.Write(w); err != nil {
		return false, err
	}

	return drift.HasDrift(results), nil
}
