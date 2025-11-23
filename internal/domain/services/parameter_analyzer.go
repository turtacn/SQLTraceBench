package services

import (
	"strconv"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type ParamType int

const (
	ParamTypeUnknown ParamType = iota
	ParamTypeInt
	ParamTypeString
	ParamTypeDatetime
)

type ParameterStats struct {
	ParamName   string
	Type        ParamType
	ValueCounts map[interface{}]int
	TotalCount  int
}

type ParameterAnalyzer struct {
	MaxCardinality int // Max unique values to track per parameter
	detector       *HotspotDetector
}

func NewParameterAnalyzer() *ParameterAnalyzer {
	return &ParameterAnalyzer{
		MaxCardinality: 10000,
		detector:       NewHotspotDetector(),
	}
}

// Analyze processes traces and returns statistical models for each parameter.
func (a *ParameterAnalyzer) Analyze(traces []models.SQLTrace) map[string]*models.ParameterModel {
	statsMap := make(map[string]*ParameterStats)
	limit := a.MaxCardinality
	if limit <= 0 {
		limit = 10000
	}

	// 1. Accumulate Statistics
	for _, trace := range traces {
		if trace.Parameters == nil {
			continue
		}
		for paramName, value := range trace.Parameters {
			stats, ok := statsMap[paramName]
			if !ok {
				stats = &ParameterStats{
					ParamName:   paramName,
					Type:        ParamTypeUnknown,
					ValueCounts: make(map[interface{}]int),
				}
				statsMap[paramName] = stats
			}

			// Type inference logic
			if stats.Type == ParamTypeUnknown {
				stats.Type = inferType(value)
			} else if stats.Type != ParamTypeString {
				currentType := inferType(value)
				if currentType != stats.Type {
					// Fallback to string if mixed types
					stats.Type = ParamTypeString
				}
			}

			// Count value
			if len(stats.ValueCounts) < limit {
				stats.ValueCounts[value]++
			} else {
				if _, exists := stats.ValueCounts[value]; exists {
					stats.ValueCounts[value]++
				}
				// else: ignore tail
			}
			stats.TotalCount++
		}
	}

	// 2. Convert to ParameterModel using HotspotDetector
	result := make(map[string]*models.ParameterModel)
	for paramName, stats := range statsMap {
		model := a.detector.DetectDistribution(stats)
		result[paramName] = model
	}

	return result
}

func inferType(value interface{}) ParamType {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return ParamTypeInt
	case time.Time:
		return ParamTypeDatetime
	case string:
		// Try parsing as int
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			return ParamTypeInt
		}
		// Try parsing as datetime
		// RFC3339, MySQL format, etc.
		if _, err := time.Parse(time.RFC3339, v); err == nil {
			return ParamTypeDatetime
		}
		if _, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return ParamTypeDatetime
		}
		if _, err := time.Parse("2006-01-02", v); err == nil {
			return ParamTypeDatetime
		}
		return ParamTypeString
	default:
		return ParamTypeString // Default fallback
	}
}
