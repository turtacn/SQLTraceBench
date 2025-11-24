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

func TestDefaultService_ValidateTrace_Success(t *testing.T) {
	// This test mocks the ValidateTrace flow.
    // Since ValidateTrace takes paths to trace files, we need to create dummy trace files.

    origTrace := models.SQLTrace{Query: "SELECT * FROM t1", Timestamp: time.Now()}
    genTrace := models.SQLTrace{Query: "SELECT * FROM t1", Timestamp: time.Now()}

    origPath := createTempTraceFile(t, []models.SQLTrace{origTrace})
    genPath := createTempTraceFile(t, []models.SQLTrace{genTrace})
    reportDir := os.TempDir()
    defer os.Remove(origPath)
    defer os.Remove(genPath)

	service := NewService()

	req := ValidationRequest{
        OriginalPath: origPath,
        GeneratedPath: genPath,
        ReportDir: reportDir,
        KSThreshold: 0.05,
    }

	report, err := service.ValidateTrace(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, report)
}

func createTempTraceFile(t *testing.T, traces []models.SQLTrace) string {
    t.Helper()
    file, err := os.CreateTemp("", "traces-*.json")
    require.NoError(t, err)
    defer file.Close()

    for _, tr := range traces {
        data, _ := json.Marshal(tr)
        file.Write(data)
        file.WriteString("\n")
    }
    return file.Name()
}
