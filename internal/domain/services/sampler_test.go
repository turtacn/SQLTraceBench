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

	// Create a distribution with 10 values
	dist := models.NewValueDistribution()
	values := []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"}
	for _, v := range values {
		// Observations don't matter for ZipfSampler, only the list of values
		dist.Values = append(dist.Values, v)
		dist.Frequencies = append(dist.Frequencies, 1)
		dist.Total++
	}

	// Sample 10,000 times
	counts := make(map[string]int)
	for i := 0; i < 10000; i++ {
		val, err := sampler.Sample(dist)
		assert.NoError(t, err)
		counts[val.(string)]++
	}

	// Verify high frequency elements appear significantly more often than low frequency elements.
	// Zipfian distribution: rank 1 (v1) should have highest frequency, rank N (v10) lowest.
	// Since ZipfSampler generates indices [0, n-1], index 0 corresponds to v1.

	t.Logf("Counts: %v", counts)

	assert.True(t, counts["v1"] > counts["v2"], "v1 should occur more than v2")
	assert.True(t, counts["v2"] > counts["v3"], "v2 should occur more than v3")
	assert.True(t, counts["v1"] > counts["v10"], "v1 should occur significantly more than v10")

	// Check if it follows power law roughly
	// v1 should be dominant.
	assert.Greater(t, counts["v1"], 2000, "v1 should have significant portion of samples")
}

func TestWeightedSampler(t *testing.T) {
	sampler := services.NewSeededWeightedRandomSampler(42)

	dist := models.NewValueDistribution()
	dist.AddObservation("A")
	dist.AddObservation("A")
	dist.AddObservation("A") // 3 A's
	dist.AddObservation("B") // 1 B

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
