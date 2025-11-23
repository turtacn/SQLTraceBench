package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func newTestDBExecutionService(t *testing.T) (*DBExecutionService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	rc := services.NewTokenBucketRateController(100, 1) // MaxConcurrency 1 to avoid race in sqlmock expectations
	return &DBExecutionService{db: db, rc: rc}, mock
}

func TestDBExecutionService_RunBench_Success(t *testing.T) {
	// Setup
	service, mock := newTestDBExecutionService(t)
	defer service.db.Close()

	workload := &models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
			{Query: "SELECT 2", Args: []interface{}{}},
		},
	}

	// Expectations
	mock.ExpectPing()

	// Query 1
	prep1 := mock.ExpectPrepare("SELECT 1")
	prep1.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	prep1.WillBeClosed()

	// Query 2
	prep2 := mock.ExpectPrepare("SELECT 2")
	prep2.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	prep2.WillBeClosed()

	// Execute
	metrics, err := service.RunBench(context.Background(), workload)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(2), metrics.QueriesExecuted)
	assert.Equal(t, int64(0), metrics.Errors)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBExecutionService_RunBench_PingFails(t *testing.T) {
	// Setup
	service, mock := newTestDBExecutionService(t)
	defer service.db.Close()

	workload := &models.BenchmarkWorkload{}

	// Expectations
	mock.ExpectPing().WillReturnError(assert.AnError)

	// Execute
	_, err := service.RunBench(context.Background(), workload)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping database")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBExecutionService_RunBench_ExecFails(t *testing.T) {
	// Setup
	service, mock := newTestDBExecutionService(t)
	defer service.db.Close()

	workload := &models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}

	// Expectations
	mock.ExpectPing()
	prep := mock.ExpectPrepare("SELECT 1")
	prep.ExpectExec().WillReturnError(assert.AnError)
	prep.WillBeClosed()

	// Execute
	metrics, err := service.RunBench(context.Background(), workload)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(1), metrics.QueriesExecuted)
	assert.Equal(t, int64(1), metrics.Errors)
	assert.NoError(t, mock.ExpectationsWereMet())
}
