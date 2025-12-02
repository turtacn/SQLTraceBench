package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// Service is the application service for the execution phase.
type Service interface {
	RunBenchmark(ctx context.Context, workload *models.BenchmarkWorkload, cfg ExecutionConfig) (*models.BenchmarkResult, error)
}

// DefaultService is the default implementation of the execution service.
type DefaultService struct {
	registry *plugin_registry.Registry
}

// NewService creates a new DefaultService.
func NewService(registry *plugin_registry.Registry) Service {
	return &DefaultService{registry: registry}
}

// ExecutionConfig holds the configuration for a benchmark run.
type ExecutionConfig struct {
	TargetDB    string
	TargetQPS   int
	Concurrency int
}

// RunBenchmark runs the benchmark.
func (s *DefaultService) RunBenchmark(ctx context.Context, workload *models.BenchmarkWorkload, cfg ExecutionConfig) (*models.BenchmarkResult, error) {
	plugin, ok := s.registry.Get(cfg.TargetDB)
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", cfg.TargetDB)
	}

	limiter := services.NewTokenBucketRateController(cfg.TargetQPS, cfg.Concurrency)
	results := make(chan models.QueryExecutionResult, 10000)
	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, q := range workload.Queries {
				limiter.Acquire(ctx)

				start := time.Now()
				// Convert args to string
				var args []string
				for _, arg := range q.Args {
					args = append(args, fmt.Sprintf("%v", arg))
				}
				_, err := plugin.ExecuteQuery(ctx, &proto.ExecuteQueryRequest{Sql: q.Query, Args: args})
				duration := time.Since(start)

				results <- models.QueryExecutionResult{Duration: duration, Error: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// This is a placeholder for a proper metrics aggregation service.
	var latencies []time.Duration
	for res := range results {
		if res.Error == nil {
			latencies = append(latencies, res.Duration)
		}
	}

	totalDuration := time.Since(startTime)
	qps := float64(len(latencies)) / totalDuration.Seconds()

	return &models.BenchmarkResult{
		Latencies: latencies,
		QPS:       qps,
	}, nil
}