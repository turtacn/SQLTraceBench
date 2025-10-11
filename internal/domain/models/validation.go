package models

import "time"

// PerformanceMetrics holds the key performance indicators from a benchmark run.
type PerformanceMetrics struct {
	// QueriesExecuted is the total number of queries run during the benchmark.
	QueriesExecuted int64
	// Duration is the total time taken for the benchmark to complete.
	Duration time.Duration
}

// QPS calculates the average queries per second.
// It returns 0 if the duration is zero to avoid division by zero errors.
func (pm *PerformanceMetrics) QPS() float64 {
	if pm.Duration.Seconds() == 0 {
		return 0
	}
	return float64(pm.QueriesExecuted) / pm.Duration.Seconds()
}

// ValidationReport contains the results of a comparison between two benchmark runs.
type ValidationReport struct {
	// BaseQPS is the QPS of the base (or control) benchmark run.
	BaseQPS float64
	// CandidateQPS is the QPS of the candidate (or test) benchmark run.
	CandidateQPS float64
	// DiffQPS is the difference in QPS between the candidate and base runs.
	DiffQPS float64
	// Threshold is the acceptable performance degradation threshold.
	Threshold float64
	// Pass indicates whether the candidate's performance is within the acceptable threshold.
	Pass bool
}