package checkpoint

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Write prints a human-readable summary of a checkpoint State to w.
func Write(w io.Writer, s State) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintf(tw, "Checkpoint recorded:\t%s\n", s.RecordedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(tw, "Total services:\t%d\n", len(s.Results))

	drifted := 0
	for _, r := range s.Results {
		if r.Drifted {
			drifted++
		}
	}
	fmt.Fprintf(tw, "Drifted services:\t%d\n", drifted)

	if drifted == 0 {
		return
	}
	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "SERVICE\tSTATUS")
	for _, r := range s.Results {
		if r.Drifted {
			fmt.Fprintf(tw, "%s\tDRIFTED\n", r.Service)
		}
	}
}
