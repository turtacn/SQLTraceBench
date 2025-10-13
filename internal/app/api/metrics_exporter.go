package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

var (
	queriesExecuted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sqltracebench_queries_executed_total",
			Help: "Total number of queries executed.",
		},
		[]string{"executor", "target"},
	)
	queryLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sqltracebench_query_latency_seconds",
			Help:    "Latency of queries in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"executor", "target"},
	)
)

func init() {
	prometheus.MustRegister(queriesExecuted)
	prometheus.MustRegister(queryLatency)
}

// MetricsExporter is a service that exports benchmark metrics in Prometheus format.
type MetricsExporter struct{}

// NewMetricsExporter creates a new MetricsExporter.
func NewMetricsExporter() *MetricsExporter {
	return &MetricsExporter{}
}

// RecordMetrics records the metrics from a benchmark run.
func (e *MetricsExporter) RecordMetrics(executor, target string, metrics *models.PerformanceMetrics) {
	labels := prometheus.Labels{"executor": executor, "target": target}
	queriesExecuted.With(labels).Add(float64(metrics.QueriesExecuted))
	for _, latency := range metrics.Latencies {
		queryLatency.With(labels).Observe(latency.Seconds())
	}
}

// Handler returns a Gin handler that serves the metrics endpoint.
func (e *MetricsExporter) Handler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}