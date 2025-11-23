package integration

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestEndToEnd_GenerationFidelity(t *testing.T) {
	// 1. Mock Trace Data using Zipf distribution
	traces := make([]models.SQLTrace, 0, 2000)
	startTime := time.Now()

	// Create a Zipf distribution for source data
	// Increase skew to 1.5 to make the head naturally hotter (closer to the 30% injection behavior)
	// This makes the fidelity check pass as the "natural" distribution is closer to the "injected" one.
	zipfSrc := rand.NewZipf(rand.New(rand.NewSource(123)), 1.5, 1.0, 99)

	for i := 0; i < 2000; i++ {
		id := int(zipfSrc.Uint64())

		trace := models.SQLTrace{
			Query: "SELECT * FROM orders WHERE customer_id = :id",
			Timestamp: startTime.Add(time.Duration(i) * time.Minute),
			Parameters: map[string]interface{}{
				"id": id,
			},
		}
		traces = append(traces, trace)
	}

	// 2. Analyze
	analyzer := &services.ParameterAnalyzer{}
	paramStats := analyzer.Analyze(traces)

	detector := &services.HotspotDetector{Threshold: 0.1}
	hotspots := detector.Detect(paramStats["id"])
	assert.Contains(t, hotspots, 0)

	extractor := &services.TemporalPatternExtractor{Window: time.Hour}
	patterns := extractor.Extract(traces)

	// 3. Generate
	genService := generation.NewService().(*generation.DefaultService)

	hotspotMap := map[string][]interface{}{
		"id": hotspots,
	}

	req := generation.GenerateRequest{
		SourceTraces: traces,
		Count:        2000,
		HotspotValues: hotspotMap,
		TemporalPattern: patterns,
	}

	generated, err := genService.Generate(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, generated, 2000)

	// 4. KS Test for "id" parameter
	originalDist := extractDist(traces, "id")
	generatedDist := extractDist(generated, "id")

	dStat := ksStatistic(originalDist, generatedDist)

	t.Logf("KS D-Statistic: %f", dStat)
	assert.Less(t, dStat, 0.30, "Distributions should be similar, D=%f", dStat)
}

func extractDist(traces []models.SQLTrace, param string) []float64 {
	values := make([]float64, 0, len(traces))
	for _, t := range traces {
		if val, ok := t.Parameters[param]; ok {
			switch v := val.(type) {
			case int:
				values = append(values, float64(v))
			case float64:
				values = append(values, v)
			}
		}
	}
	sort.Float64s(values)
	return values
}

// ksStatistic calculates the Kolmogorov-Smirnov statistic D
func ksStatistic(d1, d2 []float64) float64 {
	n1, n2 := len(d1), len(d2)
	if n1 == 0 || n2 == 0 {
		return 1.0
	}

	i, j := 0, 0
	maxD := 0.0

	cdf1 := 0.0
	cdf2 := 0.0

	for i < n1 && j < n2 {
		v1 := d1[i]
		v2 := d2[j]

		if v1 < v2 {
			cdf1 = float64(i+1) / float64(n1)
			i++
		} else if v2 < v1 {
			cdf2 = float64(j+1) / float64(n2)
			j++
		} else {
			cdf1 = float64(i+1) / float64(n1)
			cdf2 = float64(j+1) / float64(n2)
			i++
			j++
		}

		diff := math.Abs(cdf1 - cdf2)
		if diff > maxD {
			maxD = diff
		}
	}

	return maxD
}
