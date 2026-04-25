// Package window implements a sliding time-window aggregator for drift results.
//
// It is useful for computing short-term drift rates and suppressing noise from
// transient failures by examining only observations within a recent time range.
//
// Example usage:
//
//	w, err := window.New(5 * time.Minute)
//	if err != nil { ... }
//	w.Add(results)
//	fmt.Println(w.DriftedCount())
package window
