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

	zipfSrc := rand.NewZipf(rand.New(rand.NewSource(123)), 1.5, 1.0, 99)

	for i := 0; i < 2000; i++ {
		id := int(zipfSrc.Uint64())

		trace := models.SQLTrace{
			Query:     "SELECT * FROM orders WHERE customer_id = :customer_id",
			Timestamp: startTime.Add(time.Duration(i) * time.Minute),
			Parameters: map[string]interface{}{
				":customer_id": id,
			},
		}
		traces = append(traces, trace)
	}

	// 2. Analyze
	analyzer := services.NewParameterAnalyzer()
	paramModels := analyzer.Analyze(traces)

	model, ok := paramModels[":customer_id"]
	require.True(t, ok, "parameter model for :customer_id should exist")
	require.NotNil(t, model)

	assert.Equal(t, models.DistZipfian, model.DistType)
	assert.Equal(t, 0, model.TopValues[0])

	// 3. Generate
	genService := generation.NewService()

	req := generation.GenerateRequest{
		SourceTraces: traces,
		Count:        2000,
	}

	workload, err := genService.GenerateWorkload(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, workload.Queries, 2000)

	// 4. KS Test for "customer_id" parameter
	originalDist := extractDistFromTraces(traces, ":customer_id")
	generatedDist := extractIntsFromWorkloadForQuery(workload, "SELECT * FROM orders WHERE customer_id = :customer_id")

	dStat := ksStatistic(originalDist, generatedDist)

	t.Logf("KS D-Statistic: %f", dStat)
	assert.Less(t, dStat, 0.40, "Distributions should be similar, D=%f", dStat)
}

func extractIntsFromWorkloadForQuery(workload *models.BenchmarkWorkload, query string) []float64 {
	values := make([]float64, 0)
	for _, q := range workload.Queries {
		if q.Query == query {
			for _, arg := range q.Args {
				if val, ok := arg.(int); ok {
					values = append(values, float64(val))
				} else if val, ok := arg.(float64); ok {
					values = append(values, val)
				}
			}
		}
	}
	sort.Float64s(values)
	return values
}

func extractDistFromTraces(traces []models.SQLTrace, param string) []float64 {
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

func extractDistFromWorkload(workload *models.BenchmarkWorkload, argIndex int) []float64 {
	values := make([]float64, 0, len(workload.Queries))
	for _, q := range workload.Queries {
		if len(q.Args) > argIndex {
			if val, ok := q.Args[argIndex].(int); ok {
				values = append(values, float64(val))
			} else if val, ok := q.Args[argIndex].(float64); ok {
				values = append(values, val)
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
	cdf1, cdf2 := 0.0, 0.0

	for i < n1 && j < n2 {
		v1, v2 := d1[i], d2[j]
		if v1 <= v2 {
			cdf1 = float64(i+1) / float64(n1)
			i++
		}
		if v2 <= v1 {
			cdf2 = float64(j+1) / float64(n2)
			j++
		}
		diff := math.Abs(cdf1 - cdf2)
		if diff > maxD {
			maxD = diff
		}
	}
	return maxD
}
