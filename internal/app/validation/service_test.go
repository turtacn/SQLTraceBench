package validation

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

func createTempMetricsFile(t *testing.T, metrics models.PerformanceMetrics) string {
	t.Helper()
	file, err := os.CreateTemp("", "metrics-*.json")
	require.NoError(t, err)
	defer file.Close()

	err = json.NewEncoder(file).Encode(metrics)
	require.NoError(t, err)

	return file.Name()
}

func TestDefaultService_Validate_Success(t *testing.T) {
	// Create temporary metrics files.
	baseMetrics := models.PerformanceMetrics{QueriesExecuted: 100, Errors: 5, P50: 50, Duration: 10 * time.Second}
	candMetrics := models.PerformanceMetrics{QueriesExecuted: 100, Errors: 10, P50: 60, Duration: 12 * time.Second}
	basePath := createTempMetricsFile(t, baseMetrics)
	candPath := createTempMetricsFile(t, candMetrics)
	outputPath := "/tmp/report.json"
	defer os.Remove(basePath)
	defer os.Remove(candPath)
	defer os.Remove(outputPath)

	// Create the service.
	service := NewService()

	// Validate the metrics.
	report, err := service.Validate(context.Background(), basePath, candPath, outputPath, 0.1)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, report)
	assert.False(t, report.Result.Pass)
}

func TestDefaultService_Validate_BaseFileNotFound(t *testing.T) {
	// Create the service.
	service := NewService()

	// Validate with a non-existent base file.
	_, err := service.Validate(context.Background(), "non-existent.json", "cand.json", "/tmp/report.json", 0.1)

	// Assert the error.
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultService_Validate_CandFileNotFound(t *testing.T) {
	// Create a temporary base metrics file.
	baseMetrics := models.PerformanceMetrics{QueriesExecuted: 100, Errors: 5, P50: 50, Duration: 10 * time.Second}
	basePath := createTempMetricsFile(t, baseMetrics)
	defer os.Remove(basePath)

	// Create the service.
	service := NewService()

	// Validate with a non-existent candidate file.
	_, err := service.Validate(context.Background(), basePath, "non-existent.json", "/tmp/report.json", 0.1)

	// Assert the error.
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}