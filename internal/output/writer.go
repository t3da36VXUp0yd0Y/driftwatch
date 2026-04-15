// Package output provides destination-agnostic writing for drift reports.
package output

import (
	"fmt"
	"io"
	"os"
)

// Destination represents where output should be written.
type Destination int

const (
	Stdout Destination = iota
	File
)

// Options configures the output writer.
type Options struct {
	Dest     Destination
	FilePath string
}

// Writer wraps an io.WriteCloser with its destination metadata.
type Writer struct {
	w       io.WriteCloser
	managed bool // true if we opened the file and must close it
}

// New creates a Writer based on the provided Options.
// If Dest is File, the file is created or truncated.
// Callers must call Close when done.
func New(opts Options) (*Writer, error) {
	switch opts.Dest {
	case Stdout:
		return &Writer{w: os.Stdout, managed: false}, nil
	case File:
		if opts.FilePath == "" {
			return nil, fmt.Errorf("output: file path must not be empty when destination is File")
		}
		f, err := os.Create(opts.FilePath)
		if err != nil {
			return nil, fmt.Errorf("output: failed to create file %q: %w", opts.FilePath, err)
		}
		return &Writer{w: f, managed: true}, nil
	default:
		return nil, fmt.Errorf("output: unknown destination %d", opts.Dest)
	}
}

// Write implements io.Writer.
func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

// Close closes the underlying writer only if it was opened by New.
func (w *Writer) Close() error {
	if w.managed {
		return w.w.Close()
	}
	return nil
}
