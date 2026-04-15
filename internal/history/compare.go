package history

import "github.com/driftwatch/internal/drift"

// Diff describes a change in drift status for a single service between
// two successive history entries.
type Diff struct {
	Service      string
	WasDrifted   bool
	IsDrifted    bool
	StatusChanged bool
}

// Compare returns the per-service drift status changes between two entries.
// Services present only in one entry are included with zero-value booleans
// for the missing side.
func Compare(prev, curr Entry) []Diff {
	prevMap := indexResults(prev.Results)
	currMap := indexResults(curr.Results)

	seen := make(map[string]struct{})
	var diffs []Diff

	for svc, currResult := range currMap {
		seen[svc] = struct{}{}
		prevResult, hadPrev := prevMap[svc]
		d := Diff{
			Service:   svc,
			IsDrifted: currResult.Drifted,
		}
		if hadPrev {
			d.WasDrifted = prevResult.Drifted
		}
		d.StatusChanged = d.WasDrifted != d.IsDrifted
		diffs = append(diffs, d)
	}

	for svc, prevResult := range prevMap {
		if _, ok := seen[svc]; ok {
			continue
		}
		diffs = append(diffs, Diff{
			Service:       svc,
			WasDrifted:    prevResult.Drifted,
			IsDrifted:     false,
			StatusChanged: prevResult.Drifted,
		})
	}

	return diffs
}

func indexResults(results []drift.Result) map[string]drift.Result {
	m := make(map[string]drift.Result, len(results))
	for _, r := range results {
		m[r.Service] = r
	}
	return m
}
