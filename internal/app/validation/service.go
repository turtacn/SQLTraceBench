package validation

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"gonum.org/v1/gonum/stat"
)

// Service is the interface for the validation service.
type Service interface {
	ValidateBenchmarks(ctx context.Context, base, cand *models.BenchmarkResult) (*models.ValidationReport, error)
}

// DefaultService is the default implementation of the validation service.
type DefaultService struct{}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{}
}

// ValidateBenchmarks performs statistical validation between a base and candidate benchmark result.
func (s *DefaultService) ValidateBenchmarks(ctx context.Context, base, cand *models.BenchmarkResult) (*models.ValidationReport, error) {
	report := &models.ValidationReport{}

	// 1. QPS Comparison
	if base.QPS > 0 {
		report.QPSDeviation = math.Abs(base.QPS-cand.QPS) / base.QPS
	}

	// 2. Latency Profile
	baseLatencies := base.Latencies
	candLatencies := cand.Latencies

	if len(baseLatencies) > 0 && len(candLatencies) > 0 {
		// Sort the latencies before calculating quantiles
		sort.Slice(baseLatencies, func(i, j int) bool { return baseLatencies[i] < baseLatencies[j] })
		sort.Slice(candLatencies, func(i, j int) bool { return candLatencies[i] < candLatencies[j] })
		report.LatencyP99Diff = stat.Quantile(0.99, stat.Empirical, float64s(candLatencies), nil) -
			stat.Quantile(0.99, stat.Empirical, float64s(baseLatencies), nil)
	}

	// 3. Score
	if report.QPSDeviation < 0.1 {
		report.Status = "PASS"
	} else {
		report.Status = "WARN"
	}

	return report, nil
}

// Helper functions
// These are no longer needed here if we only compare BenchmarkResult, but keeping them for now.

func float64s(d []time.Duration) []float64 {
	f := make([]float64, len(d))
	for i, v := range d {
		f[i] = float64(v)
	}
	return f
}
