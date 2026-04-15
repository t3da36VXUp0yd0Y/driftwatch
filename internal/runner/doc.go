// Package runner ties together config loading, Docker inspection, drift
// detection, and report generation into a single orchestration layer.
//
// Typical usage:
//
//	cfg, _ := config.Load("driftwatch.yaml")
//	client, _ := docker.NewClient()
//	r := runner.New(cfg, client)
//	rep, _ := r.Run(ctx, report.FormatText)
//	rep.Write(os.Stdout)
package runner
