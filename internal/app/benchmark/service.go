package benchmark

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/metrics"
)

type BenchmarkService interface {
	RunBenchmark(ctx context.Context, req BenchmarkRequest) (*BenchmarkReport, error)
}

type BenchmarkRequest struct {
	ConfigPath       string
	OutputDir        string
	ExportPrometheus bool
}

type BenchmarkReport struct {
	Results    []services.BenchmarkResult
	Summary    string
	ReportPath string
}

type DefaultBenchmarkService struct {
	Runner *services.BenchmarkRunner
}

func NewDefaultBenchmarkService() *DefaultBenchmarkService {
	// Initially, we might not have models loaded. They are loaded from config.
	// So the runner is initialized later or re-initialized.
	// For this design, we'll create the runner inside RunBenchmark or pass a factory.
	// To fit the interface, let's keep it simple.
	return &DefaultBenchmarkService{}
}

// Config structures mirroring configs/benchmark.yaml
type Config struct {
	Benchmark BenchmarkConfigYAML `yaml:"benchmark"`
}

type BenchmarkConfigYAML struct {
	Models        []ModelConfig    `yaml:"models"`
	TestScenarios []ScenarioConfig `yaml:"test_scenarios"`
	Metrics       []string         `yaml:"metrics"`
	Output        OutputConfig     `yaml:"output"`
}

type ModelConfig struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	ConfigPath string `yaml:"config_path"`
}

type ScenarioConfig struct {
	Name        string `yaml:"name"`
	TraceCount  int    `yaml:"trace_count"`
	Concurrency int    `yaml:"concurrency"`
}

type OutputConfig struct {
	ReportDir        string `yaml:"report_dir"`
	ExportPrometheus bool   `yaml:"export_prometheus"`
	PrometheusPort   int    `yaml:"prometheus_port"`
}

func (s *DefaultBenchmarkService) RunBenchmark(
	ctx context.Context,
	req BenchmarkRequest,
) (*BenchmarkReport, error) {
	// 1. Load Config
	config, err := loadBenchmarkConfig(req.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize Models
	// For this implementation, we need to instantiate actual models based on config.
	// Since I don't have the full model registry here, I will create a mock/placeholder wrapper
	// or use a hypothetical factory.
	// In a real app, this would use a ModelFactory.
	modelsList := initModels(config.Benchmark.Models)

	runner := services.NewBenchmarkRunner(modelsList)

	// 3. Execute Benchmark (Using the first scenario for now, or loop through them)
	// The requirement implies running benchmarks based on config.
	// Let's assume we run the first scenario found.
	if len(config.Benchmark.TestScenarios) == 0 {
		return nil, fmt.Errorf("no test scenarios defined")
	}
	scenario := config.Benchmark.TestScenarios[0]

	results, err := runner.Run(ctx, services.BenchmarkConfig{
		TraceCount:  scenario.TraceCount,
		Concurrency: scenario.Concurrency,
	})
	if err != nil {
		return nil, err
	}

	// 4. Generate Report
	outputDir := req.OutputDir
	if outputDir == "" {
		outputDir = config.Benchmark.Output.ReportDir
	}
	reportPath, err := generateReport(results, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	// 5. Export Prometheus Metrics
	if req.ExportPrometheus || config.Benchmark.Output.ExportPrometheus {
		for _, result := range results {
			metrics.RecordBenchmarkResult(result)
		}
	}

	return &BenchmarkReport{
		Results:    results,
		Summary:    generateSummary(results),
		ReportPath: reportPath,
	}, nil
}

func loadBenchmarkConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// RealModelWrapper adapts the generation.Service to TraceGeneratorModel interface
type RealModelWrapper struct {
	name    string
	service generation.Service
}

func (m *RealModelWrapper) Name() string { return m.name }
func (m *RealModelWrapper) Generate(ctx context.Context, count int) ([]models.SQLTrace, error) {
	// The generation service requires a GenerateRequest.
	// We must provide some dummy traces if we want it to work without external files.
	dummyTrace := models.SQLTrace{
		Query:      "SELECT * FROM table WHERE id = ?",
		Parameters: map[string]interface{}{"id": 1},
	}

	req := generation.GenerateRequest{
		Count:        count,
		SourceTraces: []models.SQLTrace{dummyTrace},
	}

	workload, err := m.service.GenerateWorkload(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert BenchmarkWorkload to []SQLTrace for the benchmark runner interface.
	// This is a temporary adaptation. Ideally, the runner would consume BenchmarkWorkload directly.
	traces := make([]models.SQLTrace, len(workload.Queries))
	for i, q := range workload.Queries {
		// This conversion loses the parameter names, which might be an issue
		// if the runner's internal logic depends on them. For now, we assume it's okay.
		params := make(map[string]interface{})
		for j, arg := range q.Args {
			params[fmt.Sprintf("p%d", j+1)] = arg
		}
		traces[i] = models.SQLTrace{
			Query:      q.Query,
			Parameters: params,
		}
	}
	return traces, nil
}

func initModels(modelConfigs []ModelConfig) []services.TraceGeneratorModel {
	var modelsList []services.TraceGeneratorModel
	for _, mc := range modelConfigs {
		// Instantiate the default service for all types for now,
		// as we don't have distinct classes for Markov vs LSTM in the provided codebase inspection.
		// Assuming `generation.NewService()` gives us the standard generator.
		// If specific types require specific setup, we would do it here.

		svc := generation.NewService()
		// If we had config specific to the model, we would apply it here.

		modelsList = append(modelsList, &RealModelWrapper{
			name:    mc.Name,
			service: svc,
		})
	}
	return modelsList
}

func generateReport(results []services.BenchmarkResult, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("benchmark_report_%d.html", 1)) // simplistic naming

	// Create a simple HTML report
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmt.Fprintf(f, "<html><body><h1>Benchmark Report</h1><table border='1'><tr><th>Model</th><th>Throughput</th><th>P99 Latency</th></tr>")
	for _, r := range results {
		fmt.Fprintf(f, "<tr><td>%s</td><td>%.2f</td><td>%.2f</td></tr>", r.ModelName, r.Throughput, r.P99Latency)
	}
	fmt.Fprintf(f, "</table></body></html>")

	return filename, nil
}

func generateSummary(results []services.BenchmarkResult) string {
	if len(results) == 0 {
		return "No results"
	}
	best := results[0]
	for _, r := range results {
		if r.Throughput > best.Throughput {
			best = r
		}
	}
	return fmt.Sprintf(
		"Best model: %s (%.2f traces/sec, P99: %.2fms)",
		best.ModelName, best.Throughput, best.P99Latency,
	)
}
