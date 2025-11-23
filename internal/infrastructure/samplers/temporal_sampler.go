package samplers

import (
	"math/rand"
	"sort"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

type TemporalSampler struct {
	Pattern    *services.TemporalPattern
	BaseTime   time.Time
	weights    []float64
	cumWeights []float64
	rand       *rand.Rand
}

func NewTemporalSampler(pattern *services.TemporalPattern, baseTime time.Time) *TemporalSampler {
	s := &TemporalSampler{
		Pattern:  pattern,
		BaseTime: baseTime,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	s.init()
	return s
}

func (s *TemporalSampler) init() {
	if s.Pattern == nil || len(s.Pattern.BinCounts) == 0 {
		return
	}

	// Determine max bin index to know size
	maxBin := 0
	totalCount := 0
	for idx, count := range s.Pattern.BinCounts {
		if idx > maxBin {
			maxBin = idx
		}
		totalCount += count
	}

	size := maxBin + 1
	s.weights = make([]float64, size)
	s.cumWeights = make([]float64, size)

	var cum float64
	for i := 0; i < size; i++ {
		count := s.Pattern.BinCounts[i]
		weight := float64(count) / float64(totalCount)
		s.weights[i] = weight
		cum += weight
		s.cumWeights[i] = cum
	}
	// Ensure last is 1.0
	if size > 0 {
		s.cumWeights[size-1] = 1.0
	}
}

func (s *TemporalSampler) SampleTimestamp() time.Time {
	if s.Pattern == nil || len(s.cumWeights) == 0 {
		return s.BaseTime
	}

	r := s.rand.Float64()
	binIndex := sort.Search(len(s.cumWeights), func(i int) bool {
		return s.cumWeights[i] >= r
	})

	// Random offset within the bin
	// Window is Duration
	windowNs := s.Pattern.Window.Nanoseconds()
	offsetNs := s.rand.Int63n(windowNs)

	return s.BaseTime.Add(time.Duration(binIndex) * s.Pattern.Window).Add(time.Duration(offsetNs))
}
