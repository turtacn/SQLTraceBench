package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestValidationE2E_FullPipeline(t *testing.T) {
	// 1. Prepare test data
	baseResult := &models.BenchmarkResult{
		Latencies: []time.Duration{
			10 * time.Millisecond,
			20 * time.Millisecond,
			30 * time.Millisecond,
		},
		QPS: 1.5,
	}

	candResult := &models.BenchmarkResult{
		Latencies: []time.Duration{
			100 * time.Millisecond,
			200 * time.Millisecond,
			300 * time.Millisecond,
		},
		QPS: 1,
	}

	// 2. Initialize Service and execute
	svc := validation.NewService()

	// Run validation
	report, err := svc.ValidateBenchmarks(context.Background(), baseResult, candResult)
	require.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "WARN", report.Status)
}
