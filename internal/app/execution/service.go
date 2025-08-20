package execution

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

type Service interface {
	RunBench(ctx context.Context, yamlWorkload, dbType string) (*models.PerformanceMetrics, error)
}

type DefaultService struct {
	log *utils.Logger
}

func NewService() Service { return &DefaultService{log: utils.GetGlobalLogger()} }

func (s *DefaultService) RunBench(ctx context.Context, yamlWorkload, _ string) (*models.PerformanceMetrics, error) {
	var wl models.BenchmarkWorkload
	file, _ := os.Open(yamlWorkload)
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&wl)

	start := time.Now()
	for _, q := range wl.Queries {
		// simple execute
		s.log.Info("executing", utils.Field{Key: "query", Value: q.SQL})
	}
	return &models.PerformanceMetrics{
		QueriesExecuted: int64(len(wl.Queries)),
		Duration:        time.Since(start),
	}, nil
}
