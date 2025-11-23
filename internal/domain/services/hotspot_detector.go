package services

import (
	"sort"
)

type HotspotDetector struct {
	Threshold float64 // e.g. 0.05 for Top 5%
}

type valueFrequency struct {
	Value interface{}
	Count int
}

func (d *HotspotDetector) Detect(stats *ParameterStats) []interface{} {
	if stats == nil || stats.TotalCount == 0 {
		return nil
	}

	// Flatten map to slice for sorting
	valueFreqs := make([]valueFrequency, 0, len(stats.ValueCounts))
	for v, c := range stats.ValueCounts {
		valueFreqs = append(valueFreqs, valueFrequency{Value: v, Count: c})
	}

	// Sort by frequency descending
	sort.Slice(valueFreqs, func(i, j int) bool {
		return valueFreqs[i].Count > valueFreqs[j].Count
	})

	targetFreq := float64(stats.TotalCount) * d.Threshold
	var hotspots []interface{}
	cumulative := 0

	for _, item := range valueFreqs {
		// Optimization: If a single value is very rare (e.g. 1 occurrence), it's hardly a hotspot
		// unless the total count is small.
		// Maybe we enforce a minimum count? For now, stick to pure threshold logic.

		hotspots = append(hotspots, item.Value)
		cumulative += item.Count
		if float64(cumulative) >= targetFreq {
			break
		}
	}

	// Heuristic: If we selected too many values (e.g. > 10% of total unique values),
	// it might not be a "hotspot" distribution (uniform distribution case).
	// But the requirement says "Detect hotspots using frequency threshold".
	// We'll stick to the basic logic.

	return hotspots
}
