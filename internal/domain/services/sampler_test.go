package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestWeightedRandomSampler_Sample(t *testing.T) {
	sampler := NewWeightedRandomSampler()

	// Create a distribution with known weights
	dist := &models.ValueDistribution{
		Values:      []interface{}{"A", "B", "C"},
		Frequencies: []int{10, 30, 60}, // A: 10%, B: 30%, C: 60%
		Total:       100,
	}

	// Sample 1000 times and check the distribution
	counts := make(map[interface{}]int)
	for i := 0; i < 1000; i++ {
		sample, err := sampler.Sample(dist)
		assert.NoError(t, err)
		counts[sample]++
	}

	// Check that the distribution is roughly correct
	assert.InDelta(t, 100, counts["A"], 50, "A should be sampled about 10% of the time")
	assert.InDelta(t, 300, counts["B"], 50, "B should be sampled about 30% of the time")
	assert.InDelta(t, 600, counts["C"], 50, "C should be sampled about 60% of the time")
}