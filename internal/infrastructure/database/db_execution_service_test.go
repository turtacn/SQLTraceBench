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
	rc := services.NewTokenBucketRateController(100, 10)
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
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("SELECT 2").WillReturnResult(sqlmock.NewResult(1, 1))

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
	mock.ExpectExec("SELECT 1").WillReturnError(assert.AnError)

	// Execute
	metrics, err := service.RunBench(context.Background(), workload)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, int64(1), metrics.QueriesExecuted)
	assert.Equal(t, int64(1), metrics.Errors)
	assert.NoError(t, mock.ExpectationsWereMet())
}