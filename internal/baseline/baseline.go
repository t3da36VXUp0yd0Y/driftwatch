// Package baseline provides functionality for capturing and comparing
// a known-good configuration state against current drift results.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Entry represents a single baseline record for a service.
type Entry struct {
	Service string `json:"service"`
	Image   string `json:"image"`
	CapturedAt time.Time `json:"captured_at"`
}

// Baseline holds the full set of baseline entries.
type Baseline struct {
	Entries    []Entry   `json:"entries"`
	CapturedAt time.Time `json:"captured_at"`
}

// Capture creates a new Baseline from the provided drift results,
// recording only services that are currently healthy (no drift).
func Capture(results []drift.Result) *Baseline {
	entries := make([]Entry, 0, len(results))
	for _, r := range results {
		if !r.Drifted {
			entries = append(entries, Entry{
				Service:    r.Service,
				Image:      r.RunningImage,
				CapturedAt: time.Now().UTC(),
			})
		}
	}
	return &Baseline{
		Entries:    entries,
		CapturedAt: time.Now().UTC(),
	}
}

// Save writes the baseline to the given file path as JSON.
func Save(b *Baseline, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("baseline: create directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// Load reads a Baseline from the given file path.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Baseline{}, nil
		}
		return nil, fmt.Errorf("baseline: open file: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode: %w", err)
	}
	return &b, nil
}
