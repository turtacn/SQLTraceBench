package integration

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/app/benchmark"
)

// MockModel for integration test
type MockModel struct {
	name string
	delay time.Duration
}

func (m *MockModel) Name() string { return m.name }
func (m *MockModel) Generate(ctx context.Context, count int) ([]models.SQLTrace, error) {
	time.Sleep(m.delay) // Simulate generation time
	return make([]models.SQLTrace, count), nil
}

func TestMultiModelComparison(t *testing.T) {
	// 1. Configure 3 models
	modelsList := []services.TraceGeneratorModel{
		&MockModel{name: "markov_v1", delay: 1 * time.Millisecond},
		&MockModel{name: "lstm_v2", delay: 2 * time.Millisecond},
		&MockModel{name: "transformer", delay: 3 * time.Millisecond},
	}

	// 2. Run benchmark
	runner := services.NewBenchmarkRunner(modelsList)
	results, err := runner.Run(context.Background(), services.BenchmarkConfig{
		TraceCount:  100,
		Concurrency: 5,
		Timeout:     10 * time.Second,
	})

	assert.NoError(t, err)

	// 3. Verify results
	assert.Len(t, results, 3)
	for _, result := range results {
		assert.Greater(t, result.Throughput, 0.0)
		assert.Less(t, result.P99Latency, 5000.0)
	}

	// 4. Verify report generation (via service)
    // We reuse the service logic here or just test the service separately.
    // Let's test the service flow with a temporary config.

    // Create a temporary config file
    configContent := `
benchmark:
  models:
    - name: "markov_v1"
  test_scenarios:
    - name: "test_load"
      trace_count: 10
      concurrency: 2
  output:
    report_dir: "./test_reports"
    export_prometheus: false
`
    configFile := "test_benchmark_config.yaml"
    ioutil.WriteFile(configFile, []byte(configContent), 0644)
    defer func() {
         // Cleanup
    }()

    svc := benchmark.NewDefaultBenchmarkService()
    report, err := svc.RunBenchmark(context.Background(), benchmark.BenchmarkRequest{
        ConfigPath: configFile,
        OutputDir: "./test_reports",
    })

    assert.NoError(t, err)
    assert.FileExists(t, report.ReportPath)

    content, _ := ioutil.ReadFile(report.ReportPath)
    assert.Contains(t, string(content), "markov_v1")
}
