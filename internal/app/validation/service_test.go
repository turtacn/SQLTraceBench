package validation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestValidationService(t *testing.T) {
	service := NewService()

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

	report, err := service.ValidateBenchmarks(context.Background(), baseResult, candResult)
	assert.NoError(t, err)
	assert.Equal(t, "WARN", report.Status)
	assert.InDelta(t, 0.333, report.QPSDeviation, 0.01)
}
