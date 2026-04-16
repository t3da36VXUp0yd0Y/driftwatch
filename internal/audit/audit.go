// Package audit records drift check events to an append-only audit log.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Service   string         `json:"service"`
	Drifted   bool           `json:"drifted"`
	Details   []drift.Result `json:"details"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// New returns a Logger that appends to the file at path.
func New(path string) (*Logger, error) {
	if path == "" {
		return nil, fmt.Errorf("audit: path must not be empty")
	}
	return &Logger{path: path}, nil
}

// Record appends an audit entry for each result.
func (l *Logger) Record(results []drift.Result) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, r := range results {
		e := Entry{
			Timestamp: time.Now().UTC(),
			Service:   r.Service,
			Drifted:   r.Drifted,
			Details:   []drift.Result{r},
		}
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("audit: encode entry: %w", err)
		}
	}
	return nil
}

// Load reads all audit entries from the log file.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
