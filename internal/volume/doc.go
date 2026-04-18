// Package volume checks for drift in container volume mount declarations.
// It compares the expected mounts defined in configuration against those
// reported by the Docker client for each running service.
package volume
