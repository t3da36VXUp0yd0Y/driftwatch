// Package schedule provides a Scheduler that periodically invokes a Runner
// to detect configuration drift between deployed services and their declared
// infrastructure state.
//
// The Scheduler runs drift checks at a fixed interval. Each check delegates
// to a Runner implementation, which compares live infrastructure state against
// the declared configuration and reports any discrepancies to the provided
// writer.
//
// Usage:
//
//	s, err := schedule.New(runner, 30*time.Second, os.Stdout)
//	if err != nil { ... }
//	if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) { ... }
//
// Start blocks until the provided context is cancelled. To run the scheduler
// in the background, invoke it in a goroutine:
//
//	go func() {
//		if err := s.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
//			log.Printf("scheduler stopped unexpectedly: %v", err)
//		}
//	}()
package schedule
