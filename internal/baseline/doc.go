// Package baseline captures a known-good state of running services
// and allows future drift results to be compared against that baseline.
//
// A baseline is typically captured after a successful deployment and
// stored as a JSON file on disk. Subsequent runs can load the baseline
// to determine whether drift represents a new regression or a
// previously accepted state.
package baseline
