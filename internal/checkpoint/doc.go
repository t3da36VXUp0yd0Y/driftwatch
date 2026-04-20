// Package checkpoint persists drift results between runs so that driftwatch
// can determine whether the infrastructure state has changed since the last
// execution. Use Save to write a checkpoint and Load to restore it. The
// Changed helper compares a previous checkpoint against fresh results.
package checkpoint
