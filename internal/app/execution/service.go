package execution

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// Service is the interface for the execution service.
type Service interface {
	RunBench(ctx context.Context, workloadPath string, qps, concurrency int, slowThreshold time.Duration) (*models.PerformanceMetrics, error)
}

// DefaultService is the default implementation of the execution service.
type DefaultService struct{}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{}
}

// RunBench executes a benchmark workload from a file.
func (s *DefaultService) RunBench(ctx context.Context, workloadPath string, qps, concurrency int, slowThreshold time.Duration) (*models.PerformanceMetrics, error) {
	file, err := os.Open(workloadPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wl models.BenchmarkWorkload
	if err := json.NewDecoder(file).Decode(&wl); err != nil {
		return nil, err
	}

	rc := services.NewTokenBucketRateController(qps, concurrency)
	execSvc := services.NewExecutionService(rc, slowThreshold)

	metrics, err := execSvc.RunBench(ctx, &wl)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}