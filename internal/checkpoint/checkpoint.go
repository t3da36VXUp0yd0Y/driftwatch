// Package checkpoint persists the last-known drift state so consecutive
// runs can detect whether drift has appeared, resolved, or remained unchanged.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// State represents a persisted checkpoint entry.
type State struct {
	RecordedAt time.Time            `json:"recorded_at"`
	Results    []drift.Result       `json:"results"`
}

// Save writes the current drift results to the given file path.
func Save(path string, results []drift.Result) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	s := State{
		RecordedAt: time.Now().UTC(),
		Results:    results,
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

// Load reads a previously saved checkpoint from disk.
// Returns an empty State (no error) when the file does not exist.
func Load(path string) (State, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	defer f.Close()
	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return State{}, err
	}
	return s, nil
}

// Changed reports whether the set of drifted service names differs between
// the previous checkpoint and the current results.
func Changed(prev State, current []drift.Result) bool {
	prevSet := make(map[string]bool, len(prev.Results))
	for _, r := range prev.Results {
		if r.Drifted {
			prevSet[r.Service] = true
		}
	}
	currSet := make(map[string]bool, len(current))
	for _, r := range current {
		if r.Drifted {
			currSet[r.Service] = true
		}
	}
	if len(prevSet) != len(currSet) {
		return true
	}
	for k := range currSet {
		if !prevSet[k] {
			return true
		}
	}
	return false
}
