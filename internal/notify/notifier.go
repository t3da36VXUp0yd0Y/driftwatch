// Package notify provides notification hooks for drift detection results.
package notify

import (
	"fmt"
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// Level controls which results trigger a notification.
type Level int

const (
	// LevelAll sends notifications for every result.
	LevelAll Level = iota
	// LevelDriftOnly sends notifications only when drift is detected.
	LevelDriftOnly
)

// Notifier sends drift notifications to a destination.
type Notifier struct {
	w     io.Writer
	level Level
}

// New returns a Notifier that writes to w and filters by level.
func New(w io.Writer, level Level) *Notifier {
	return &Notifier{w: w, level: level}
}

// Notify writes a notification for each result that passes the level filter.
// It returns the number of notifications written and any write error.
func (n *Notifier) Notify(results []drift.Result) (int, error) {
	count := 0
	for _, r := range results {
		if n.level == LevelDriftOnly && !r.Drifted {
			continue
		}
		if err := n.write(r); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (n *Notifier) write(r drift.Result) error {
	status := "OK"
	if r.Drifted {
		status = "DRIFT"
	}
	_, err := fmt.Fprintf(n.w, "[%s] service=%s expected=%s actual=%s reason=%s\n",
		status, r.Service, r.Expected, r.Actual, r.Reason)
	return err
}
