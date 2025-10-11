// Package services contains the interfaces for the application's core services.
package services

import (
	"context"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ExecutionService is responsible for running a benchmark workload and collecting performance metrics.
// This is a simulated implementation that uses `time.Sleep` to mimic database latency.
type ExecutionService struct{}

// NewExecutionService creates a new ExecutionService.
func NewExecutionService() *ExecutionService {
	return &ExecutionService{}
}

// RunBench executes a benchmark workload and returns the performance metrics.
// It iterates through the queries in the workload, simulating execution with a short sleep.
// The context can be used to cancel the benchmark run prematurely.
func (s *ExecutionService) RunBench(ctx context.Context, wl models.BenchmarkWorkload) (*models.PerformanceMetrics, error) {
	start := time.Now()

	for _, query := range wl.Queries {
		select {
		case <-ctx.Done():
			// The context was canceled, so stop the benchmark.
			return nil, ctx.Err()
		default:
			// Simulate the execution of the query.
			_ = query // a read to avoid "unused variable" compiler error
			time.Sleep(1 * time.Millisecond)
		}
	}

	duration := time.Since(start)
	metrics := &models.PerformanceMetrics{
		QueriesExecuted: int64(len(wl.Queries)),
		Duration:        duration,
	}

	return metrics, nil
}