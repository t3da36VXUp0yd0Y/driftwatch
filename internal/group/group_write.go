package group

import (
	"fmt"
	"io"
)

// Write renders a summary of each group to w.
func Write(w io.Writer, groups []Group) error {
	for _, g := range groups {
		total := len(g.Results)
		drifted := g.DriftedCount()
		status := "ok"
		if drifted > 0 {
			status = "DRIFT"
		}
		_, err := fmt.Fprintf(w, "[%s] group=%-20s total=%-4d drifted=%-4d status=%s\n",
			status, g.Name, total, drifted, status)
		if err != nil {
			return err
		}
	}
	return nil
}
