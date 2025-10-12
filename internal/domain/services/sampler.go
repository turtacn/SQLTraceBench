package services

import (
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// Sampler is the interface for parameter value samplers.
// It defines the contract for any sampling algorithm.
type Sampler interface {
	// Sample selects a value from the distribution based on the sampling strategy.
	Sample(dist *models.ValueDistribution) (interface{}, error)
}

// WeightedRandomSampler selects a value from a distribution based on its observed frequency.
// This creates a more realistic workload where common values appear more often.
type WeightedRandomSampler struct {
	rand *rand.Rand
}

// NewWeightedRandomSampler creates a new sampler with a random seed.
func NewWeightedRandomSampler() *WeightedRandomSampler {
	return &WeightedRandomSampler{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewSeededWeightedRandomSampler creates a new sampler with a specific seed for deterministic behavior.
func NewSeededWeightedRandomSampler(seed int64) *WeightedRandomSampler {
	return &WeightedRandomSampler{
		rand: rand.New(rand.NewSource(seed)),
	}
}

// Sample performs weighted random sampling on a given value distribution.
// It calculates a random number and iterates through the values, subtracting their
// probability mass from the random number until it drops below zero.
func (s *WeightedRandomSampler) Sample(dist *models.ValueDistribution) (interface{}, error) {
	if dist.Total == 0 {
		return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
	}

	r := s.rand.Float64()
	var cumulativeProb float64 = 0.0

	for i, value := range dist.Values {
		probability := float64(dist.Frequencies[i]) / float64(dist.Total)
		cumulativeProb += probability
		if r < cumulativeProb {
			return value, nil
		}
	}

	// Fallback in case of floating point inaccuracies
	return dist.Values[len(dist.Values)-1], nil
}