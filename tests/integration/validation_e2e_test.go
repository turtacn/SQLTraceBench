package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"context"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func createTestTraceFile(t *testing.T, path string, count int, paramOffset float64) {
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	for i := 0; i < count; i++ {
		trace := models.SQLTrace{
			Query:     "SELECT * FROM users WHERE id = ?",
			Timestamp: time.Now(),
			Parameters: map[string]interface{}{
				"id": float64(i) + paramOffset,
			},
		}
		bytes, _ := json.Marshal(trace)
		file.Write(bytes)
		file.WriteString("\n")
	}
}

func TestValidationE2E_FullPipeline(t *testing.T) {
	// 1. Prepare test data
	tempDir, err := ioutil.TempDir("", "validation_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalPath := filepath.Join(tempDir, "original.jsonl")
	generatedPath := filepath.Join(tempDir, "generated.jsonl")
	reportDir := filepath.Join(tempDir, "reports")

	createTestTraceFile(t, originalPath, 100, 0)
	createTestTraceFile(t, generatedPath, 100, 0.5) // Slight shift, should pass KS with correct threshold

	// 2. Initialize Service and execute
	svc := validation.NewService()
	req := validation.ValidationRequest{
		OriginalPath:   originalPath,
		GeneratedPath:  generatedPath,
		ReportDir:      reportDir,
		PrometheusPort: 0, // Disable for this test or use random port
		KSThreshold:    0.01,
	}

	// We use a random high port for Prometheus to avoid conflicts
	req.PrometheusPort = 19091

	// Run validation
	report, err := svc.ValidateTrace(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, report)
	assert.Greater(t, report.OverallScore, 50.0)

	// 3. Verify HTML Report
	reportPath := filepath.Join(reportDir, "validation_report.html")
	assert.FileExists(t, reportPath)

	content, err := ioutil.ReadFile(reportPath)
	require.NoError(t, err)
	htmlContent := string(content)
	assert.Contains(t, htmlContent, "Distribution Test Results")
	assert.Contains(t, htmlContent, "id") // Parameter name

	// 4. Verify Prometheus Metrics
	// Wait a bit for server to start (it's in a goroutine)
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", req.PrometheusPort))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	metricsContent := string(body)
	assert.Contains(t, metricsContent, "validation_ks_pvalue")
	assert.Contains(t, metricsContent, "validation_pass_rate")
}
