// Package pipeline provides a composable, stage-based processing pipeline
// for drift detection results.
//
// Stages implement the Stage interface and are executed in the order they
// are registered. Each stage receives the output of the previous stage,
// allowing filtering, annotation, and transformation to be composed
// cleanly without coupling individual concerns.
//
// Example usage:
//
//	p, err := pipeline.New(filterStage, maskStage, tagStage)
//	if err != nil { ... }
//	results, err := p.Run(ctx, rawResults)
package pipeline
