package services

import (
	"sync"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// MetricsRecorder is responsible for recording performance metrics during a benchmark run.
// It is designed to be thread-safe.
type MetricsRecorder struct {
	mu          sync.Mutex
	metrics     *models.PerformanceMetrics
	slowThreshold time.Duration
}

// NewMetricsRecorder creates a new metrics recorder.
// - slowThreshold: The duration after which a query is considered "slow".
func NewMetricsRecorder(slowThreshold time.Duration) *MetricsRecorder {
	return &MetricsRecorder{
		metrics: &models.PerformanceMetrics{
			Latencies: make([]time.Duration, 0),
		},
		slowThreshold: slowThreshold,
	}
}

// Record records the result of a single query execution.
// It captures the latency and whether an error occurred.
func (r *MetricsRecorder) Record(latency time.Duration, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.QueriesExecuted++
	r.metrics.Latencies = append(r.metrics.Latencies, latency)

	if err != nil {
		r.metrics.Errors++
	}

	if latency > r.slowThreshold {
		r.metrics.SlowQueries++
	}
}

// Finalize calculates the summary statistics after the benchmark run is complete.
// It should be called once at the end of the benchmark.
func (r *MetricsRecorder) Finalize(totalDuration time.Duration) *models.PerformanceMetrics {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics.Duration = totalDuration
	r.metrics.CalculatePercentiles()

	return r.metrics
}