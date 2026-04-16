// Package plugin provides a registry for named drift-check plugins.
//
// Plugins are user-supplied CheckFn functions that return drift results.
// They are registered by name and can be executed individually or all at once
// via RunAll, which merges results across all registered plugins.
package plugin
