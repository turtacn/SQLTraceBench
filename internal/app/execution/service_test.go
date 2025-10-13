package execution

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// MockExecutor is a mock implementation of the Executor interface for testing.
type MockExecutor struct {
	RunBenchFunc func(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error)
}

func (m *MockExecutor) RunBench(ctx context.Context, wl *models.BenchmarkWorkload) (*models.PerformanceMetrics, error) {
	if m.RunBenchFunc != nil {
		return m.RunBenchFunc(ctx, wl)
	}
	return &models.PerformanceMetrics{}, nil
}

func createTempWorkloadFile(t *testing.T, workload models.BenchmarkWorkload) string {
	t.Helper()
	file, err := os.CreateTemp("", "workload-*.json")
	require.NoError(t, err)
	defer file.Close()

	err = json.NewEncoder(file).Encode(workload)
	require.NoError(t, err)

	return file.Name()
}

func TestDefaultService_RunBench_Simulated(t *testing.T) {
	// Create a temporary workload file.
	workload := models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}
	workloadPath := createTempWorkloadFile(t, workload)
	defer os.Remove(workloadPath)

	// Create the service.
	service := NewService()

	// Run the benchmark with the simulated executor.
	metrics, err := service.RunBench(context.Background(), workloadPath, "simulated", "", "", 10, 1, 1*time.Second)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, metrics)
}

func TestDefaultService_RunBench_UnknownExecutor(t *testing.T) {
	// Create a temporary workload file.
	workload := models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}
	workloadPath := createTempWorkloadFile(t, workload)
	defer os.Remove(workloadPath)

	// Create the service.
	service := NewService()

	// Run the benchmark with an unknown executor.
	_, err := service.RunBench(context.Background(), workloadPath, "unknown", "", "", 10, 1, 1*time.Second)

	// Assert the error.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown executor type: unknown")
}

func TestDefaultService_RunBench_WorkloadNotFound(t *testing.T) {
	// Create the service.
	service := NewService()

	// Run the benchmark with a non-existent workload file.
	_, err := service.RunBench(context.Background(), "non-existent-file.json", "simulated", "", "", 10, 1, 1*time.Second)

	// Assert the error.
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultService_RunBench_Real_ConnectionError(t *testing.T) {
	// This test expects a connection error because it uses a dummy DSN.
	// This verifies that the "real" executor path is taken.

	// Create a temporary workload file.
	workload := models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}
	workloadPath := createTempWorkloadFile(t, workload)
	defer os.Remove(workloadPath)

	// Create the service.
	service := NewService()

	// Run the benchmark with the real executor and a bad DSN.
	// We need a registered driver for sql.Open to succeed. "mysql" is in go.mod.
	_, err := service.RunBench(context.Background(), workloadPath, "real", "mysql", "user:pass@tcp(127.0.0.1:3306)/db", 10, 1, 1*time.Second)

	// Assert the error.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping database")
}