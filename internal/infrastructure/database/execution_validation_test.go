package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestExecutionAndValidation(t *testing.T) {
	// 1. Setup Mock DB
	// Use MonitorPingsOption(true) to match ExpectPing
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	// IMPORTANT: Allow out-of-order execution because requests run concurrently
	mock.MatchExpectationsInOrder(false)

	// 2. Setup Workload
	workload := &models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{},
	}
	// Generate 50 queries
	for i := 0; i < 50; i++ {
		workload.Queries = append(workload.Queries, models.QueryWithArgs{
			Query: "SELECT * FROM users WHERE id = ?",
			Args:  []interface{}{i},
		})
	}

	// Expectation: 50 queries
	mock.ExpectPing()
	queryRegex := "SELECT \\* FROM users WHERE id = \\?"

	// Expect Prepare ONCE because we use the same query string
	// Since order doesn't matter, we can just define it.
	// Note: With MatchExpectationsInOrder(false), the Prepare might happen anytime.
	prep := mock.ExpectPrepare(queryRegex)

	// Expect 50 executions
	for i := 0; i < 50; i++ {
		// Expect concurrent execution, order doesn't matter
		prep.ExpectExec().WithArgs(i).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	prep.WillBeClosed()

	// 3. Setup Services
	// Create Rate Controller
	// MaxConcurrency=1 to reduce race conditions in test environment, though code handles it.
	// Even with MaxConcurrency=1, Go routines start instantly, so concurrency happens.
	rc := services.NewTokenBucketRateController(1000, 1)

	// Manually inject DB into DBExecutionService
	svc := &DBExecutionService{
		db: db,
		rc: rc,
	}
	// Configure connection pool
	db.SetMaxOpenConns(rc.MaxConcurrency())
	db.SetMaxIdleConns(rc.MaxConcurrency())

	// 4. Run Workload
	ctx := context.Background()
	start := time.Now()
	metrics, err := svc.RunBench(ctx, workload)
	duration := time.Since(start)

	require.NoError(t, err)

	if metrics.Errors > 0 {
		t.Logf("Metrics Errors: %d", metrics.Errors)
	}
	assert.Equal(t, int64(0), metrics.Errors)
	assert.Equal(t, int64(50), metrics.QueriesExecuted)

	// Check duration roughly. 50 queries at 1000 QPS = 50ms.
	// Allow buffer.
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(45))

	// 5. Validation Logic (Testing ValidationService)
	// Create dummy baseline metrics
	baseMetrics := &models.PerformanceMetrics{
		QueriesExecuted: 50,
		Duration:        duration,
		P99:             10 * time.Millisecond,
	}
	baseMetrics.CalculatePercentiles()

	// Validation Service
	validationSvc := services.NewValidationService()
	metadata := &models.ReportMetadata{
		Threshold: 0.05, // 5%
	}

	// Create a temp file for report
	tmpFile := "test_report"
	defer os.Remove(tmpFile + ".json")
	defer os.Remove(tmpFile + ".html")

	report, err := validationSvc.ValidateAndReport(baseMetrics, metrics, metadata, tmpFile)
	require.NoError(t, err)
	assert.NotNil(t, report)

	// Check if report files exist
	_, err = os.Stat(tmpFile + ".json")
	assert.NoError(t, err)
	_, err = os.Stat(tmpFile + ".html")
	assert.NoError(t, err)

	assert.True(t, report.Result.Pass)
	assert.NoError(t, mock.ExpectationsWereMet())
}
