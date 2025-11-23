package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestHotspotDetector(t *testing.T) {
	detector := NewHotspotDetector()

	t.Run("TestUniformDistribution", func(t *testing.T) {
		// [1, 2, 3, 4, 5] each 10 times
		stats := &ParameterStats{
			ParamName:   "uniform_param",
			Type:        ParamTypeInt,
			ValueCounts: make(map[interface{}]int),
		}
		for i := 1; i <= 5; i++ {
			stats.ValueCounts[i] = 10
			stats.TotalCount += 10
		}

		model := detector.DetectDistribution(stats)

		assert.Equal(t, models.DistUniform, model.DistType)
		assert.Less(t, model.HotspotRatio, 0.3) // 1/5 = 0.2.
	})

	t.Run("TestZipfDistribution", func(t *testing.T) {
		// Simulate Zipf-like data
		// 1 -> 1000
		// 2 -> 500
		// 3 -> 333
		// ...
		stats := &ParameterStats{
			ParamName:   "zipf_param",
			Type:        ParamTypeInt,
			ValueCounts: make(map[interface{}]int),
		}

		for i := 1; i <= 100; i++ {
			count := int(1000.0 / float64(i)) // Zipf with s=1
			stats.ValueCounts[i] = count
			stats.TotalCount += count
		}

		model := detector.DetectDistribution(stats)

		assert.Equal(t, models.DistZipfian, model.DistType)
		assert.Greater(t, model.ZipfS, 0.8)
		assert.Less(t, model.ZipfS, 1.3) // Should be around 1.0

		fmt.Printf("Detected Zipf S: %f, Ratio: %f\n", model.ZipfS, model.HotspotRatio)
	})

	t.Run("TestTopValues", func(t *testing.T) {
		stats := &ParameterStats{
			ParamName:   "top_param",
			Type:        ParamTypeString,
			ValueCounts: map[interface{}]int{
				"A": 100,
				"B": 50,
				"C": 10,
			},
			TotalCount: 160,
		}

		model := detector.DetectDistribution(stats)

		assert.Equal(t, 3, len(model.TopValues))
		assert.Equal(t, "A", model.TopValues[0])
		assert.Equal(t, 100, model.TopFrequencies[0])
		assert.Equal(t, 0.625, model.HotspotRatio) // 100/160 = 0.625
		assert.Equal(t, models.DistZipfian, model.DistType) // Ratio > 0.4
	})
}
