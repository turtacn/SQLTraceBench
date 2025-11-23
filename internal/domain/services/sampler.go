package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/math/distributions"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// Sampler is the interface for parameter value samplers.
type Sampler interface {
	// Sample selects a value from the distribution based on the sampling strategy.
	Sample(dist *models.ParameterModel) (interface{}, error)
}

// WeightedRandomSampler selects a value from a distribution based on its observed frequency.
type WeightedRandomSampler struct {
	rand *rand.Rand
}

func NewWeightedRandomSampler() *WeightedRandomSampler {
	return &WeightedRandomSampler{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewSeededWeightedRandomSampler(seed int64) *WeightedRandomSampler {
	return &WeightedRandomSampler{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (s *WeightedRandomSampler) Sample(dist *models.ParameterModel) (interface{}, error) {
	if dist.Cardinality == 0 && len(dist.TopValues) == 0 {
		return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
	}

	// Use TopFrequencies for weighted random if available
	if len(dist.TopFrequencies) > 0 {
		total := 0
		for _, f := range dist.TopFrequencies {
			total += f
		}

		if total == 0 {
			// Fallback to uniform on top values
			idx := s.rand.Intn(len(dist.TopValues))
			return dist.TopValues[idx], nil
		}

		r := s.rand.Float64()
		var cumulativeProb float64 = 0.0

		for i, value := range dist.TopValues {
			probability := float64(dist.TopFrequencies[i]) / float64(total)
			cumulativeProb += probability
			if r < cumulativeProb {
				return value, nil
			}
		}
		// Fallback
		return dist.TopValues[len(dist.TopValues)-1], nil
	}

	// If no frequencies, use Uniform
	idx := s.rand.Intn(len(dist.TopValues))
	return dist.TopValues[idx], nil
}

// ZipfSampler selects a value from a distribution based on a Zipfian distribution.
// It prioritizes the stored ZipfS parameter if available.
type ZipfSampler struct {
	rand *rand.Rand
	s    float64 // Default Skewness
	v    float64 // Default v
}

func NewZipfSampler(s float64) *ZipfSampler {
	return &ZipfSampler{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		s:    s,
		v:    1.0,
	}
}

func NewSeededZipfSampler(seed int64, s float64) *ZipfSampler {
	return &ZipfSampler{
		rand: rand.New(rand.NewSource(seed)),
		s:    s,
		v:    1.0,
	}
}

func (z *ZipfSampler) Sample(dist *models.ParameterModel) (interface{}, error) {
	n := dist.Cardinality
	if n == 0 {
		// Try using TopValues length
		n = len(dist.TopValues)
	}
	if n == 0 {
		return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
	}

	// If there's only one value, just return it.
	if n == 1 && len(dist.TopValues) > 0 {
		return dist.TopValues[0], nil
	}

	// Use model's S if available and valid (>1.001 to be safe for Zipf)
	s := z.s
	if dist.ZipfS > 1.001 {
		s = dist.ZipfS
	}

	// Ensure S is > 1 for math/rand.Zipf
	if s <= 1.0 {
		s = 1.0001
	}

	// math/rand.Zipf requires imax (upper bound inclusive).
	imax := uint64(n - 1)

	gen := distributions.NewZipfGeneratorWithRand(z.rand, s, z.v, imax)
	idx := gen.Uint64()

	if idx >= uint64(n) {
		idx = 0
	}

	// Map index to value
	if int(idx) < len(dist.TopValues) {
		return dist.TopValues[int(idx)], nil
	}

	// Synthetic Tail Generation
	// If index is outside stored TopValues, generate a deterministic synthetic value.

	if dist.DataType == "INT" {
		// Return the index itself as the value (assuming ID-like behavior)
		return int(idx), nil
	}

	// For STRING, generate a synthetic string
	// "tail_<idx>" ensures uniqueness and follows the distribution rank
	return fmt.Sprintf("tail_%d", idx), nil
}
