// Package watch provides file-based config watching for driftwatch,
// triggering a callback whenever the config file changes on disk.
package watch

import (
	"crypto/sha256"
	"io"
	"os"
	"time"
)

// Watcher polls a file for changes and calls OnChange when a change is detected.
type Watcher struct {
	path     string
	interval time.Duration
	lastHash [sha256.Size]byte
	OnChange func()
	stop     chan struct{}
}

// New creates a Watcher for the given file path and poll interval.
// OnChange is called when the file content changes.
func New(path string, interval time.Duration, onChange func()) (*Watcher, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if onChange == nil {
		return nil, ErrNilCallback
	}
	w := &Watcher{
		path:     path,
		interval: interval,
		OnChange: onChange,
		stop:     make(chan struct{}),
	}
	h, err := w.hashFile()
	if err != nil {
		return nil, err
	}
	w.lastHash = h
	return w, nil
}

// Start begins polling. It blocks until Stop is called.
func (w *Watcher) Start() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h, err := w.hashFile()
			if err != nil {
				continue
			}
			if h != w.lastHash {
				w.lastHash = h
				w.OnChange()
			}
		case <-w.stop:
			return
		}
	}
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) hashFile() ([sha256.Size]byte, error) {
	f, err := os.Open(w.path)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [sha256.Size]byte{}, err
	}
	var out [sha256.Size]byte
	copy(out[:], h.Sum(nil))
	return out, nil
}
