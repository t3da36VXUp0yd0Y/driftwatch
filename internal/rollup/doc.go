// Package rollup provides grouping utilities that aggregate drift.Result
// slices into named buckets. This is useful when services follow a naming
// convention (e.g. "team-service") and operators want a per-team view of
// configuration drift without scanning every individual result.
package rollup
