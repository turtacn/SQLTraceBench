package services

import (
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/math/distributions"
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
	if len(dist.Values) > 0 {
		return dist.Values[len(dist.Values)-1], nil
	}
	return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
}

// ZipfSampler selects a value from a distribution based on a Zipfian distribution.
// It ignores the observed frequencies and instead imposes a Zipfian distribution
// on the set of observed values (assuming they are ranked by importance or simply by index).
type ZipfSampler struct {
	rand *rand.Rand
	s    float64 // Skewness parameter (s > 1)
	v    float64 // Parameter v (v >= 1), typically 1
}

// NewZipfSampler creates a new ZipfSampler.
// s is the skewness parameter. Larger values mean more skew (more "hotspot").
func NewZipfSampler(s float64) *ZipfSampler {
	return &ZipfSampler{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		s:    s,
		v:    1.0,
	}
}

// NewSeededZipfSampler creates a new ZipfSampler with a specific seed.
func NewSeededZipfSampler(seed int64, s float64) *ZipfSampler {
	return &ZipfSampler{
		rand: rand.New(rand.NewSource(seed)),
		s:    s,
		v:    1.0,
	}
}

// Sample performs Zipfian sampling on the *indices* of the given value distribution.
// It maps the generated Zipf index to the value at that index in the Values slice.
func (z *ZipfSampler) Sample(dist *models.ValueDistribution) (interface{}, error) {
	n := len(dist.Values)
	if n == 0 {
		return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
	}

	// If there's only one value, just return it.
	if n == 1 {
		return dist.Values[0], nil
	}

	// math/rand.Zipf requires imax (upper bound inclusive).
	// So imax should be n - 1.
	imax := uint64(n - 1)

	// Create a Zipf generator for this specific distribution size.
	// Since creating a Zipf generator can be expensive if done repeatedly for large N,
	// ideally we would cache it, but for now we create it on the fly or use a lighter approach.
	// Given we are simulating hotspots, 'n' (number of unique parameter values) might not be huge.

	// We use our wrapper which uses math/rand.Zipf
	gen := distributions.NewZipfGeneratorWithRand(z.rand, z.s, z.v, imax)

	idx := gen.Uint64()

	// Safety check, though Zipf should guarantee [0, imax]
	if idx >= uint64(n) {
		idx = 0 // Fallback to most frequent
	}

	return dist.Values[idx], nil
}
