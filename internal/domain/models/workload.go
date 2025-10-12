package models

// QueryWithArgs represents a single query to be executed, with its parameters separated
// to allow for safe execution using prepared statements.
type QueryWithArgs struct {
	// Query is the SQL statement with '?' placeholders for parameters.
	Query string `json:"query"`
	// Args is a slice of arguments to be bound to the query's placeholders.
	Args []interface{} `json:"args"`
}

// BenchmarkWorkload represents a set of queries to be executed by the benchmark.
type BenchmarkWorkload struct {
	// Queries is a list of all the SQL queries and their arguments for the workload.
	Queries []QueryWithArgs `json:"queries"`
}