// Package models contains the core domain models for the SQLTraceBench application.
package models

import "time"

// SQLTrace represents a single SQL query event captured from a trace.
type SQLTrace struct {
	Query      string
	Timestamp  time.Time
	Latency    time.Duration
	Parameters map[string]interface{}
}

// TraceCollection holds a collection of SQLTraces.
type TraceCollection struct {
	Traces []SQLTrace
}

// Add appends a new SQLTrace to the collection.
func (tc *TraceCollection) Add(trace SQLTrace) {
	tc.Traces = append(tc.Traces, trace)
}

// Len returns the number of traces in the collection.
func (tc *TraceCollection) Len() int {
	return len(tc.Traces)
}
