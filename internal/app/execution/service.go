package execution

import (
	"context"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

// Service is the interface for the execution service.
type Service interface {
	RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error)
}

// DefaultService is the default implementation of the execution service.
// It uses a simulated execution strategy.
type DefaultService struct {
	log *utils.Logger
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{log: utils.GetGlobalLogger()}
}

// RunBench executes a benchmark workload using a simulated strategy.
func (s *DefaultService) RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error) {
	start := time.Now()
	for _, q := range wl.Queries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Simulate the execution of the query.
			s.log.Info("executing", utils.Field{Key: "query", Value: q})
			time.Sleep(1 * time.Millisecond)
		}
	}
	return &models.PerformanceMetrics{
		QueriesExecuted: int64(len(wl.Queries)),
		Duration:        time.Since(start),
	}, nil
}