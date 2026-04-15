// Package snapshot provides functionality for capturing and persisting
// point-in-time views of drift detection results.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Snapshot represents a captured state of drift detection at a point in time.
type Snapshot struct {
	CapturedAt time.Time           `json:"captured_at"`
	Results    []drift.Result      `json:"results"`
	Meta       map[string]string   `json:"meta,omitempty"`
}

// New creates a new Snapshot from the given results.
func New(results []drift.Result) Snapshot {
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Results:    results,
		Meta:       make(map[string]string),
	}
}

// Save writes the snapshot as a JSON file to the given directory.
// The filename is derived from the capture timestamp.
func Save(dir string, s Snapshot) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("snapshot: create directory: %w", err)
	}

	filename := fmt.Sprintf("snapshot-%s.json", s.CapturedAt.Format("20060102-150405"))
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("snapshot: write file: %w", err)
	}

	return path, nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read file: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return s, nil
}
