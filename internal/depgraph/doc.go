// Package depgraph provides a lightweight directed dependency graph for
// driftwatch services. It allows callers to declare which services depend on
// others and then query which services are transitively affected when drift is
// detected in one or more services.
package depgraph
