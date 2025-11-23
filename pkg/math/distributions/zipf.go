package distributions

import (
	"math/rand"
	"time"
)

// ZipfGenerator wraps math/rand.Zipf to generate Zipfian distributed random numbers.
type ZipfGenerator struct {
	r *rand.Zipf
}

// NewZipfGenerator creates a new ZipfGenerator.
// s is the skewness parameter (s > 1).
// v is the parameter v (v >= 1).
// imax is the maximum integer value (inclusive) that can be generated.
// seed is the seed for the random number generator.
func NewZipfGenerator(seed int64, s float64, v float64, imax uint64) *ZipfGenerator {
	src := rand.NewSource(seed)
	rng := rand.New(src)
	z := rand.NewZipf(rng, s, v, imax)
	return &ZipfGenerator{r: z}
}

// NewZipfGeneratorWithRand creates a new ZipfGenerator using an existing rand.Rand.
func NewZipfGeneratorWithRand(rng *rand.Rand, s float64, v float64, imax uint64) *ZipfGenerator {
	z := rand.NewZipf(rng, s, v, imax)
	return &ZipfGenerator{r: z}
}

// Next returns the next Zipfian distributed random number.
func (z *ZipfGenerator) Next() uint64 {
	return z.r.Uint64()
}

// Uint64 returns the next Zipfian distributed random number.
func (z *ZipfGenerator) Uint64() uint64 {
	return z.r.Uint64()
}

// DefaultZipfGenerator creates a ZipfGenerator with current time seed and default v=1.
func DefaultZipfGenerator(s float64, imax uint64) *ZipfGenerator {
	return NewZipfGenerator(time.Now().UnixNano(), s, 1.0, imax)
}
