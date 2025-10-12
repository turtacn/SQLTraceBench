package services

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestValidationService_ValidateAndReport(t *testing.T) {
	service := NewValidationService()

	baseMetrics := &models.PerformanceMetrics{
		QueriesExecuted: 100,
		Duration:        1 * time.Second,
	}
	candMetrics := &models.PerformanceMetrics{
		QueriesExecuted: 90,
		Duration:        1 * time.Second,
	}
	metadata := &models.ReportMetadata{
		Threshold: 0.05,
	}
	outputPath := "test_report.json"
	defer os.Remove(outputPath)

	report, err := service.ValidateAndReport(baseMetrics, candMetrics, metadata, outputPath)
	require.NoError(t, err)

	assert.False(t, report.Result.Pass)
	assert.Contains(t, report.Result.Reason, "Validation failed")
	assert.Contains(t, report.Result.Reason, "90.00")
	assert.Contains(t, report.Result.Reason, "100.00")
	assert.Contains(t, report.Result.Reason, "5.00%")
}