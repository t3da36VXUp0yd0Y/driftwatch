// Package filter provides post-detection filtering of drift results.
//
// After the drift detector produces a full set of results, callers may
// use Apply to narrow the output by drift status or by a specific set of
// service names before handing the results to the report package.
//
// Example usage:
//
//	results := detector.Detect(declared, running)
//	filtered := filter.Apply(results, filter.Options{
//		OnlyDrifted:  true,
//		ServiceNames: []string{"api", "worker"},
//	})
package filter
