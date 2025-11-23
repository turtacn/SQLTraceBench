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
}

func (a *ParameterAnalyzer) Analyze(traces []models.SQLTrace) map[string]*ParameterStats {
	result := make(map[string]*ParameterStats)
	limit := a.MaxCardinality
	if limit <= 0 {
		limit = 10000 // Default limit
	}

	// Pre-process to identify all parameters
	for _, trace := range traces {
		if trace.Parameters == nil {
			continue
		}
		for paramName, value := range trace.Parameters {
			stats, ok := result[paramName]
			if !ok {
				stats = &ParameterStats{
					ParamName:   paramName,
					Type:        ParamTypeUnknown,
					ValueCounts: make(map[interface{}]int),
				}
				result[paramName] = stats
			}

			// Type inference logic
			if stats.Type == ParamTypeUnknown {
				stats.Type = inferType(value)
			} else if stats.Type != ParamTypeString {
				currentType := inferType(value)
				if currentType != stats.Type {
					stats.Type = ParamTypeString
				}
			}

			// Count value with cardinality protection
			// If map is full and value is new, skip or replace?
			// Simple approach: skip new values if limit reached.
			// Better approach (SpaceSaving) is complex.
			// Given this is a benchmark tool, we likely want "Heavy Hitters".
			// If we just drop new ones, we might miss late-arriving heavy hitters.
			// But implementing a full sketch is out of scope.
			// We'll stick to: if len < limit OR value exists, increment.
			if len(stats.ValueCounts) < limit {
				stats.ValueCounts[value]++
			} else {
				if _, exists := stats.ValueCounts[value]; exists {
					stats.ValueCounts[value]++
				}
				// else: ignore tail value
			}
			stats.TotalCount++
		}
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
