package models

import (
	"time"
)

type ValidationReport struct {
	OriginalStats  PerformanceMetrics
	SyntheticStats PerformanceMetrics
	DeviationQPS   float64 // percentage
	Passed         bool
	GeneratedAt    time.Time
}

type PerformanceMetrics struct {
	QueriesExecuted int64
	Duration        time.Duration
	AvgQPS          float64
	MinLatency      float64
	MaxLatency      float64
}

func (s *PerformanceMetrics) QPS() float64 {
	if s.Duration <= 0 {
		return 0
	}
	return float64(s.QueriesExecuted) / s.Duration.Seconds()
}

func (r *ValidationReport) Compare(original, synthetic PerformanceMetrics) {
	r.OriginalStats = original
	r.SyntheticStats = synthetic
	delta := synthetic.QPS() - original.QPS()
	if original.QPS() != 0 {
		r.DeviationQPS = (delta / original.QPS()) * 100
	}
	r.Passed = abs(r.DeviationQPS) <= 50 // relaxed threshold for MVP
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

//Personal.AI order the ending
