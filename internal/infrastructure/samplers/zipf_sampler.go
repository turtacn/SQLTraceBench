package samplers

import (
	"math/rand"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/math/distributions"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// ZipfSampler selects a value from a distribution based on a Zipfian distribution.
type ZipfSampler struct {
	rand          *rand.Rand
	s             float64 // Skewness parameter (s > 1)
	v             float64 // Parameter v (v >= 1), typically 1
	HotspotValues []interface{}
	HotspotProb   float64 // Probability to inject a hotspot value [0.0, 1.0]
}

// NewZipfSampler creates a new ZipfSampler.
func NewZipfSampler(s float64) *ZipfSampler {
	return &ZipfSampler{
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
		s:           s,
		v:           1.0,
		HotspotProb: 0.3, // Default from P2 requirements
	}
}

// Sample selects a value from the distribution.
// It supports hotspot injection.
func (z *ZipfSampler) Sample(dist *models.ValueDistribution) (interface{}, error) {
	// Hotspot injection logic
	if len(z.HotspotValues) > 0 {
		if z.rand.Float64() < z.HotspotProb {
			return z.HotspotValues[z.rand.Intn(len(z.HotspotValues))], nil
		}
	}

	n := len(dist.Values)
	if n == 0 {
		return nil, types.NewError(types.ErrInvalidInput, "cannot sample from an empty distribution")
	}

	if n == 1 {
		return dist.Values[0], nil
	}

	imax := uint64(n - 1)
	gen := distributions.NewZipfGeneratorWithRand(z.rand, z.s, z.v, imax)
	idx := gen.Uint64()

	if idx >= uint64(n) {
		idx = 0
	}

	return dist.Values[idx], nil
}
