// Package history records drift check results over time, enabling
// trend analysis and change detection between successive runs.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a single recorded drift check run.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Results   []drift.Result `json:"results"`
}

// Store manages persistence of drift history entries.
type Store struct {
	path string
}

// New creates a new Store that reads and writes to the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Append adds a new entry to the history file, creating it if necessary.
func (s *Store) Append(results []drift.Result) error {
	entries, err := s.Load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("history: load existing: %w", err)
	}

	entries = append(entries, Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	})

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("history: write file: %w", err)
	}
	return nil
}

// Load reads all history entries from disk.
// Returns an empty slice if the file does not exist.
func (s *Store) Load() ([]Entry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("history: read file: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return entries, nil
}

// Latest returns the most recent entry, and false if no entries exist.
func (s *Store) Latest() (Entry, bool) {
	entries, err := s.Load()
	if err != nil || len(entries) == 0 {
		return Entry{}, false
	}
	return entries[len(entries)-1], true
}
