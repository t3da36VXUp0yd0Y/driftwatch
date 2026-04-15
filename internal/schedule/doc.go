// Package schedule provides a Scheduler that periodically invokes a Runner
// to detect configuration drift between deployed services and their declared
// infrastructure state.
//
// Usage:
//
//	s, err := schedule.New(runner, 30*time.Second, os.Stdout)
//	if err != nil { ... }
//	if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) { ... }
package schedule
