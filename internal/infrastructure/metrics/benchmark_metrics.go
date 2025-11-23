package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

var (
	GenerationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "benchmark_generation_duration_seconds",
			Help:    "Trace generation duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"model"},
	)

	GenerationThroughput = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "benchmark_generation_throughput",
			Help: "Traces generated per second",
		},
		[]string{"model"},
	)

	ValidationScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "benchmark_validation_score",
			Help: "Validation score (0-1)",
		},
		[]string{"model"},
	)

	MemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "benchmark_memory_usage_mb",
			Help: "Memory usage in MB",
		},
		[]string{"model"},
	)
)

func init() {
	prometheus.MustRegister(
		GenerationDuration,
		GenerationThroughput,
		ValidationScore,
		MemoryUsage,
	)
}

func RecordBenchmarkResult(result services.BenchmarkResult) {
	GenerationThroughput.WithLabelValues(result.ModelName).Set(result.Throughput)
	ValidationScore.WithLabelValues(result.ModelName).Set(result.ValidationScore)
	MemoryUsage.WithLabelValues(result.ModelName).Set(result.MemoryUsageMB)
}
