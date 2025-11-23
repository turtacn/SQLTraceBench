package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestZipfSampler(t *testing.T) {
	// Initialize Zipf distribution with Skew=1.1
	sampler := services.NewSeededZipfSampler(42, 1.1)

	// Create a parameter model with 10 values
	values := []interface{}{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"}
	dist := &models.ParameterModel{
		TopValues:   values,
		Cardinality: 10,
	}

	// Sample 10,000 times
	counts := make(map[string]int)
	for i := 0; i < 10000; i++ {
		val, err := sampler.Sample(dist)
		assert.NoError(t, err)
		counts[val.(string)]++
	}

	t.Logf("Counts: %v", counts)

	assert.True(t, counts["v1"] > counts["v2"], "v1 should occur more than v2")
	assert.True(t, counts["v2"] > counts["v3"], "v2 should occur more than v3")
	assert.True(t, counts["v1"] > counts["v10"], "v1 should occur significantly more than v10")

	assert.Greater(t, counts["v1"], 2000, "v1 should have significant portion of samples")
}

func TestWeightedSampler(t *testing.T) {
	sampler := services.NewSeededWeightedRandomSampler(42)

	dist := &models.ParameterModel{
		TopValues:      []interface{}{"A", "B"},
		TopFrequencies: []int{3, 1},
		Cardinality:    2,
	}

	counts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		val, err := sampler.Sample(dist)
		assert.NoError(t, err)
		counts[val.(string)]++
	}

	t.Logf("Weighted Counts: %v", counts)
	// A should be ~750, B ~250
	assert.Greater(t, counts["A"], counts["B"])
	assert.InDelta(t, 750, counts["A"], 100)
}
