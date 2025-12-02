package models

import (
	"sort"
	"time"
)

// PerformanceMetrics holds the key performance indicators from a benchmark run.
type PerformanceMetrics struct {
	// QueriesExecuted is the total number of queries run during the benchmark.
	QueriesExecuted int64 `json:"queries_executed"`
	// Errors is the total number of errors that occurred during the benchmark.
	Errors int64 `json:"errors"`
	// Duration is the total time taken for the benchmark to complete.
	Duration time.Duration `json:"duration"`
	// Latencies is a slice of all the individual query latencies.
	Latencies []time.Duration `json:"-"` // Exclude from JSON report for brevity
	// P50 is the 50th percentile latency.
	P50 time.Duration `json:"p50"`
	// P90 is the 90th percentile latency.
	P90 time.Duration `json:"p90"`
	// P99 is the 99th percentile latency.
	P99 time.Duration `json:"p99"`
	// SlowQueries is the number of queries that exceeded the slow query threshold.
	SlowQueries int64 `json:"slow_queries"`
}

// QPS calculates the average queries per second.
func (pm *PerformanceMetrics) QPS() float64 {
	if pm.Duration.Seconds() == 0 {
		return 0
	}
	return float64(pm.QueriesExecuted) / pm.Duration.Seconds()
}

// ErrorRate calculates the percentage of queries that resulted in an error.
func (pm *PerformanceMetrics) ErrorRate() float64 {
	if pm.QueriesExecuted == 0 {
		return 0
	}
	return float64(pm.Errors) / float64(pm.QueriesExecuted)
}

// CalculatePercentiles computes the P50, P90, and P99 latencies from the collected latency data.
func (pm *PerformanceMetrics) CalculatePercentiles() {
	if len(pm.Latencies) == 0 {
		return
	}
	sort.Slice(pm.Latencies, func(i, j int) bool {
		return pm.Latencies[i] < pm.Latencies[j]
	})
	pm.P50 = pm.Latencies[len(pm.Latencies)/2]
	pm.P90 = pm.Latencies[int(float64(len(pm.Latencies))*0.9)]
	pm.P99 = pm.Latencies[int(float64(len(pm.Latencies))*0.99)]
}

// Report is the top-level structure for the benchmark report.
// It includes metadata and the validation results.
type Report struct {
	Version   string             `json:"version"`
	Timestamp time.Time          `json:"timestamp"`
	Metadata  *ReportMetadata    `json:"metadata"`
	Result    *ValidationResult `json:"result"`
}

// ReportMetadata contains contextual information about the benchmark run.
type ReportMetadata struct {
	BaseTarget      string          `json:"base_target"`
	CandidateTarget string          `json:"candidate_target"`
	Threshold       float64         `json:"threshold"`
	TemplateSummary []*SQLTemplate  `json:"template_summary"`
}

// ValidationResult contains the comparison of the two benchmark runs.
type ValidationResult struct {
	BaseMetrics      *PerformanceMetrics `json:"base_metrics"`
	CandidateMetrics *PerformanceMetrics `json:"candidate_metrics"`
	Pass             bool                `json:"pass"`
	Reason           string              `json:"reason"`
}

type BenchmarkResult struct {
	Latencies []time.Duration
	QPS       float64
}

type ValidationReport struct {
	Status         string
	QPSDeviation   float64
	LatencyP99Diff float64
}

type QueryExecutionResult struct {
	Duration time.Duration
	Error    error
}