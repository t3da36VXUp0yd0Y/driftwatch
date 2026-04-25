// Package suppress allows operators to silence known, accepted configuration
// drift so that recurring expected deviations do not pollute reports or trigger
// alerts. Rules can be scoped to a service name and an optional reason
// substring, and may carry an expiry time so that temporary suppressions are
// automatically lifted.
package suppress
