package reporters

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

var (
	ksPValue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validation_ks_pvalue",
			Help: "KS test p-value for parameter distributions",
		},
		[]string{"parameter"},
	)

	passRate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "validation_pass_rate",
			Help: "Overall validation pass rate (0-1)",
		},
	)
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(ksPValue, passRate)

	// Register handler once to avoid panic on re-registration
	http.Handle("/metrics", promhttp.Handler())
}

// PrometheusReporter exports validation metrics to Prometheus.
type PrometheusReporter struct {
	Port int
}

// NewPrometheusReporter creates a new PrometheusReporter.
func NewPrometheusReporter(port int) *PrometheusReporter {
	return &PrometheusReporter{
		Port: port,
	}
}

// ExportMetrics updates Prometheus metrics based on the validation report and starts the HTTP server.
func (r *PrometheusReporter) ExportMetrics(report *services.ValidationReport) error {
	// Update metrics
	passRate.Set(report.OverallScore / 100.0)

	for _, result := range report.DistributionTests {
		if paramName, ok := result.Details["parameter"].(string); ok {
			ksPValue.WithLabelValues(paramName).Set(result.PValue)
		}
	}

	// Start HTTP server
	// Handler is registered in init(), so we just start the server.
	return http.ListenAndServe(fmt.Sprintf(":%d", r.Port), nil)
}
