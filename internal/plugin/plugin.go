// Package plugin provides a simple registry for named drift-check plugins.
package plugin

import (
	"fmt"
	"sort"
	"sync"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// CheckFn is a function that performs a custom drift check and returns results.
type CheckFn func() ([]drift.Result, error)

// Registry holds named plugin check functions.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]CheckFn
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{plugins: make(map[string]CheckFn)}
}

// Register adds a named plugin to the registry.
// Returns an error if the name is already registered.
func (r *Registry) Register(name string, fn CheckFn) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %q already registered", name)
	}
	r.plugins[name] = fn
	return nil
}

// Names returns a sorted list of registered plugin names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.plugins))
	for n := range r.plugins {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Run executes the named plugin and returns its results.
func (r *Registry) Run(name string) ([]drift.Result, error) {
	r.mu.RLock()
	fn, ok := r.plugins[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugin %q not found", name)
	}
	return fn()
}

// RunAll executes all registered plugins and merges results.
func (r *Registry) RunAll() ([]drift.Result, error) {
	names := r.Names()
	var all []drift.Result
	for _, name := range names {
		res, err := r.Run(name)
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", name, err)
		}
		all = append(all, res...)
	}
	return all, nil
}
