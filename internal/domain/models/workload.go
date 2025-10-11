package models

// BenchmarkWorkload represents a set of queries to be executed by the benchmark.
// It is generated from a collection of SQLTemplates and a set of parameters.
type BenchmarkWorkload struct {
	// Queries is a list of all the SQL queries to be executed in the workload.
	Queries []string
}