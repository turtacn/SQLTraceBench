package services

import (
	"math"
	"sort"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type HotspotDetector struct {
	Threshold float64 // e.g. 0.05 for Top 5%
}

type valueFrequency struct {
	Value interface{}
	Count int
}

// DetectDistribution analyzes the frequency map and determines the distribution type and parameters.
func (d *HotspotDetector) DetectDistribution(stats *ParameterStats) *models.ParameterModel {
	model := &models.ParameterModel{
		ParamName:   stats.ParamName,
		Cardinality: len(stats.ValueCounts),
	}

	switch stats.Type {
	case ParamTypeInt:
		model.DataType = "INT"
	case ParamTypeDatetime:
		model.DataType = "DATETIME"
	default:
		model.DataType = "STRING"
	}

	if stats.TotalCount == 0 {
		model.DistType = models.DistUniform
		return model
	}

	// Flatten and sort
	valueFreqs := make([]valueFrequency, 0, len(stats.ValueCounts))
	for v, c := range stats.ValueCounts {
		valueFreqs = append(valueFreqs, valueFrequency{Value: v, Count: c})
	}

	// Sort by frequency descending
	sort.Slice(valueFreqs, func(i, j int) bool {
		return valueFreqs[i].Count > valueFreqs[j].Count
	})

	// Fill Top Values
	// Limit increased to 10000 to better capture distribution head and mid-tail
	limit := 10000
	if len(valueFreqs) < limit {
		limit = len(valueFreqs)
	}
	model.TopValues = make([]interface{}, limit)
	model.TopFrequencies = make([]int, limit)
	for i := 0; i < limit; i++ {
		model.TopValues[i] = valueFreqs[i].Value
		model.TopFrequencies[i] = valueFreqs[i].Count
	}

	// Calculate Hotspot Ratio
	top20Count := int(math.Ceil(float64(len(valueFreqs)) * 0.2))
	if top20Count < 1 {
		top20Count = 1
	}

	headSum := 0
	for i := 0; i < top20Count && i < len(valueFreqs); i++ {
		headSum += valueFreqs[i].Count
	}

	ratio := float64(headSum) / float64(stats.TotalCount)
	model.HotspotRatio = ratio

	// Detect Distribution Type using Threshold
	// We use the configured Threshold as a guide, or a standard heuristic.
	// Standard heuristic: If top 20% items have > 40% traffic (skewed).
	// If d.Threshold is meant to be the skewness threshold ratio?
	// Usually threshold is "top X%".
	// Let's stick to the heuristic.

	if ratio > 0.4 {
		model.DistType = models.DistZipfian
		model.ZipfS = d.estimateZipfS(valueFreqs)
		model.ZipfV = 1.0
	} else {
		model.DistType = models.DistUniform
	}

	return model
}

func (d *HotspotDetector) estimateZipfS(sortedFreqs []valueFrequency) float64 {
	if len(sortedFreqs) < 2 {
		return 1.1
	}

	n := len(sortedFreqs)
	// Limit to top 100 for regression efficiency and relevance
	if n > 100 { n = 100 }

	sumLogRank := 0.0
	sumLogFreq := 0.0
	sumLogRankLogFreq := 0.0
	sumLogRankSq := 0.0

	count := 0
	for i := 0; i < n; i++ {
		rank := float64(i + 1)
		freq := float64(sortedFreqs[i].Count)

		if freq <= 0 {
			continue
		}

		logRank := math.Log(rank)
		logFreq := math.Log(freq)

		sumLogRank += logRank
		sumLogFreq += logFreq
		sumLogRankLogFreq += logRank * logFreq
		sumLogRankSq += logRank * logRank
		count++
	}

	if count < 2 {
		return 1.1
	}

	num := float64(count)
	denominator := (num * sumLogRankSq) - (sumLogRank * sumLogRank)
	if denominator == 0 {
		return 1.1
	}

	slope := ((num * sumLogRankLogFreq) - (sumLogRank * sumLogFreq)) / denominator

	s := -slope
	if s < 0 {
		s = 0
	}

	return s
}

func NewHotspotDetector() *HotspotDetector {
	return &HotspotDetector{
		Threshold: 0.05,
	}
}
