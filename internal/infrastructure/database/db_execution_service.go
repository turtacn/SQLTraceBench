package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// DBExecutionService is an implementation of the execution service that runs benchmarks
// against a real database.
type DBExecutionService struct {
	db           *sql.DB
	rc           services.RateController
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

	// Use a local cache for this run to avoid state leakage and ensure thread safety per run.
	var (
		prepareCache sync.Map // map[string]*sql.Stmt
		prepareMu    sync.Mutex
	)

	// Ensure prepared statements are closed when the function exits, regardless of how it exits.
	defer func() {
		prepareCache.Range(func(key, value interface{}) bool {
			stmt := value.(*sql.Stmt)
			_ = stmt.Close()
			return true
		})
	}()

	// Helper to get prepared statement from local cache
	getStmt := func(query string) (*sql.Stmt, error) {
		if stmt, ok := prepareCache.Load(query); ok {
			return stmt.(*sql.Stmt), nil
		}

		prepareMu.Lock()
		defer prepareMu.Unlock()

		if stmt, ok := prepareCache.Load(query); ok {
			return stmt.(*sql.Stmt), nil
		}

		stmt, err := s.db.PrepareContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement: %w", err)
		}

		prepareCache.Store(query, stmt)
		return stmt, nil
	}

	var wg sync.WaitGroup

	for _, q := range wl.Queries {
		if err := s.rc.Acquire(ctx); err != nil {
			// If we can't acquire, we stop.
			// The already spawned goroutines will continue and finish.
			// We should probably wait for them if we want a clean shutdown,
			// but if ctx is canceled, they should also be canceled ideally if they use ctx.
			// However, here we just return error. The deferred cleanup will handle statements.
			// Note: If we return here, `wg.Wait()` below is skipped, so we might return before
			// pending queries finish. But since we return error, the metrics might be partial.
			// If we want to wait, we should do it.
			break
		}

		wg.Add(1)
		go func(query models.QueryWithArgs) {
			defer wg.Done()

			execStart := time.Now()

			stmt, err := getStmt(query.Query)
			if err != nil {
				recorder.Record(time.Since(execStart), err)
				return
			}

			_, err = stmt.ExecContext(ctx, query.Args...)
			latency := time.Since(execStart)
			recorder.Record(latency, err)
		}(q)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	return recorder.Finalize(totalDuration), ctx.Err()
}
