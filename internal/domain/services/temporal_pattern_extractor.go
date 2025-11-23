package services

import (
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type TemporalPattern struct {
	Window    time.Duration
	BinCounts map[int]int // key=time bin index, value=query count
}

type TemporalPatternExtractor struct {
	Window time.Duration
}

func (e *TemporalPatternExtractor) Extract(traces []models.SQLTrace) *TemporalPattern {
	if len(traces) == 0 {
		return &TemporalPattern{
			Window:    e.Window,
			BinCounts: make(map[int]int),
		}
	}

	minTime, maxTime := FindTimeRange(traces)

	// Avoid zero window
	if e.Window <= 0 {
		e.Window = time.Hour // Default
	}

	// Calculate number of bins
	// We use the range [minTime, maxTime]
	// Bin index = (t - minTime) / Window
	duration := maxTime.Sub(minTime)
	// Ensure at least one bin
	numBins := int(duration/e.Window) + 1
	binCounts := make(map[int]int, numBins)

	for _, trace := range traces {
		// Handle potential out of bounds or negative due to time diffs if traces are not sorted?
		// Traces might not be sorted, but minTime is absolute minimum.
		if trace.Timestamp.Before(minTime) {
			continue // Should not happen given FindTimeRange
		}
		binIndex := int(trace.Timestamp.Sub(minTime) / e.Window)
		binCounts[binIndex]++
	}

	return &TemporalPattern{
		Window:    e.Window,
		BinCounts: binCounts,
	}
}

func FindTimeRange(traces []models.SQLTrace) (minTime, maxTime time.Time) {
	if len(traces) == 0 {
		return
	}
	minTime = traces[0].Timestamp
	maxTime = traces[0].Timestamp

	for _, t := range traces {
		if t.Timestamp.Before(minTime) {
			minTime = t.Timestamp
		}
		if t.Timestamp.After(maxTime) {
			maxTime = t.Timestamp
		}
	}
	return
}
