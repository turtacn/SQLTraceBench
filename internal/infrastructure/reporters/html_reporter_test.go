package reporters_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/reporters"
)

func TestHTMLReporter_GenerateReport(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.html")

	// Create dummy validation result
	res := &models.ValidationResult{
		Pass: true,
		BaseMetrics: &models.PerformanceMetrics{
			QueriesExecuted: 1000,
			Duration:        10 * time.Second,
			P99:             50 * time.Millisecond,
			Latencies:       []time.Duration{10 * time.Millisecond, 20 * time.Millisecond},
		},
		CandidateMetrics: &models.PerformanceMetrics{
			QueriesExecuted: 1000,
			Duration:        10 * time.Second,
			P99:             55 * time.Millisecond,
			Latencies:       []time.Duration{12 * time.Millisecond, 22 * time.Millisecond},
		},
	}

	// Initialize Reporter
	reporter, err := reporters.NewHTMLReporter()
	require.NoError(t, err)

	// Execute
	err = reporter.GenerateReport(res, "mock-plugin", outputPath)
	require.NoError(t, err)

	// Verify File Exists
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	// Verify Content
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	html := string(content)

	assert.Contains(t, html, "SQLTraceBench Validation Report")
	assert.Contains(t, html, "mock-plugin")
	assert.Contains(t, html, "PASS")
	// Chart.js data
	assert.Contains(t, html, "QPS Over Time")
}
