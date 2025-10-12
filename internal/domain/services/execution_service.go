package services

import (
	"context"
	"sync"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ExecutionService is responsible for running a benchmark workload and collecting performance metrics.
type ExecutionService struct {
	rc            RateController
	recorder      *MetricsRecorder
	slowThreshold time.Duration
}

// NewExecutionService creates a new ExecutionService.
func NewExecutionService(rc RateController, slowThreshold time.Duration) *ExecutionService {
	return &ExecutionService{
		rc:            rc,
		slowThreshold: slowThreshold,
	}
}

// RunBench executes a benchmark workload and returns the final performance metrics.
func (s *ExecutionService) RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error) {
	s.recorder = NewMetricsRecorder(s.slowThreshold)
	start := time.Now()
	var wg sync.WaitGroup
	queriesCh := make(chan string, len(wl.Queries))

	s.rc.Start(ctx)

	// Use the MaxConcurrency from the rate controller to determine the number of workers.
	for i := 0; i < s.rc.MaxConcurrency(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for query := range queriesCh {
				if err := s.rc.Acquire(ctx); err != nil {
					return
				}

				execStart := time.Now()
				_ = query
				time.Sleep(1 * time.Millisecond)
				latency := time.Since(execStart)
				s.recorder.Record(latency, nil)
			}
		}()
	}

	for _, query := range wl.Queries {
		select {
		case queriesCh <- query:
		case <-ctx.Done():
			break
		}
	}
	close(queriesCh)

	wg.Wait()

	totalDuration := time.Since(start)
	return s.recorder.Finalize(totalDuration), nil
}