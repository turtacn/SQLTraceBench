package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
)

// Executor is an interface for running a benchmark workload.
// This allows for different execution strategies (e.g., simulated, real database).
type Executor interface {
	RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error)
}

// Service is the application service for the execution phase.
type Service interface {
	RunBench(
		ctx context.Context,
		workloadPath, executorType, driver, dsn string,
		qps, concurrency int,
		slowThreshold time.Duration,
	) (*models.PerformanceMetrics, error)
}

// DefaultService is the default implementation of the execution service.
type DefaultService struct{}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{}
}

// RunBench selects an executor based on the provided type and runs the benchmark.
func (s *DefaultService) RunBench(
	ctx context.Context,
	workloadPath, executorType, driver, dsn string,
	qps, concurrency int,
	slowThreshold time.Duration,
) (*models.PerformanceMetrics, error) {
	// Read the workload file.
	file, err := os.Open(workloadPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wl models.BenchmarkWorkload
	if err := json.NewDecoder(file).Decode(&wl); err != nil {
		return nil, err
	}

	// Create the rate controller.
	rc := services.NewTokenBucketRateController(qps, concurrency)

	// Select the executor.
	var executor Executor
	switch executorType {
	case "simulated":
		executor = services.NewExecutionService(rc, slowThreshold)
	case "real":
		executor, err = database.NewDBExecutionService(driver, dsn, rc)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown executor type: %s", executorType)
	}

	// Run the benchmark.
	metrics, err := executor.RunBench(ctx, &wl)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}