// Package timeout provides a simple wrapper for executing operations
// with a configurable deadline. It is used by driftwatch to bound
// Docker client calls and plugin executions so that a slow or
// unresponsive target cannot stall the drift-detection pipeline.
package timeout
