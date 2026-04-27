// Package signal computes a weighted drift signal score for each service.
//
// A signal score aggregates multiple drift check dimensions (image, env,
// port, volume, health) into a single numeric value that indicates how
// severely a service has drifted from its declared state. Higher scores
// indicate more significant or numerous drift events.
//
// Usage:
//
//	weights := signal.DefaultWeights()
//	scores  := signal.Compute(results, weights)
//	signal.Write(os.Stdout, scores)
package signal
