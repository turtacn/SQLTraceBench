package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// DBExecutionService is an implementation of the execution service that runs benchmarks
// against a real database.
type DBExecutionService struct {
	db *sql.DB
	rc services.RateController
}

// NewDBExecutionService creates a new DBExecutionService.
func NewDBExecutionService(driver, dsn string, rc services.RateController) (*DBExecutionService, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to open database connection", err)
	}

	// Configure the connection pool.
	db.SetMaxOpenConns(rc.MaxConcurrency())
	db.SetMaxIdleConns(rc.MaxConcurrency())

	return &DBExecutionService{db: db, rc: rc}, nil
}

// RunBench executes a benchmark workload against the database.
func (s *DBExecutionService) RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error) {
	// Ping the database to ensure the connection is alive.
	if err := s.db.PingContext(ctx); err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to ping database", err)
	}

	recorder := services.NewMetricsRecorder(100 * time.Millisecond) // TODO: Make slow threshold configurable.
	start := time.Now()

	s.rc.Start(ctx)
	defer s.rc.Stop()

	// In a real implementation, we would use a worker pool here.
	// For now, we'll execute queries sequentially to keep it simple.
	for _, q := range wl.Queries {
		if err := s.rc.Acquire(ctx); err != nil {
			return nil, err
		}

		execStart := time.Now()
		_, err := s.db.ExecContext(ctx, q.Query, q.Args...)
		latency := time.Since(execStart)
		recorder.Record(latency, err)
	}

	totalDuration := time.Since(start)
	return recorder.Finalize(totalDuration), nil
}