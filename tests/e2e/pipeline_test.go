package e2e

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/app/workflow"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

// MockPlugin implements plugins.Plugin interface
type MockPlugin struct {
	mock.Mock
}

func (m *MockPlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*proto.ExecuteQueryResponse), args.Error(1)
}

func (m *MockPlugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	args := m.Called(schema)
	return args.Get(0).(*models.Schema), args.Error(1)
}

func (m *MockPlugin) TranslateQuery(query string) (string, error) {
	args := m.Called(query)
	return args.String(0), args.Error(1)
}

func (m *MockPlugin) Name() string { return "mock_plugin" }
func (m *MockPlugin) Version() string { return "1.0.0" }
func (m *MockPlugin) Init(config map[string]string) error { return nil }
func (m *MockPlugin) Close() error { return nil }
func (m *MockPlugin) HealthCheck() error { return nil }

func TestFullPipeline_Mocked(t *testing.T) {
	// 1. Setup Environment
	tmpDir, err := os.MkdirTemp("", "e2e_pipeline_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create inputs
	tracePath := filepath.Join(tmpDir, "traces.jsonl")
	schemaPath := filepath.Join(tmpDir, "schema.sql")
	baselinePath := filepath.Join(tmpDir, "baseline.json")

	// Input Traces
	traceData := `{"query": "SELECT * FROM users WHERE id = :id", "timestamp": "2024-01-01T00:00:00Z", "parameters": {"id": 1}}
{"query": "SELECT * FROM users WHERE id = :id", "timestamp": "2024-01-01T00:00:01Z", "parameters": {"id": 2}}`
	err = os.WriteFile(tracePath, []byte(traceData), 0644)
	require.NoError(t, err)

	// Input Schema
	schemaData := "CREATE TABLE users (id INT PRIMARY KEY);"
	err = os.WriteFile(schemaPath, []byte(schemaData), 0644)
	require.NoError(t, err)

	// Baseline
	baseline := models.BenchmarkResult{QPS: 100, Latencies: []time.Duration{time.Millisecond}}
	baselineBytes, _ := json.Marshal(baseline)
	err = os.WriteFile(baselinePath, baselineBytes, 0644)
	require.NoError(t, err)

	// 2. Setup Services & Mocks
	mockP := new(MockPlugin)
	// Expect schema conversion
	mockP.On("ConvertSchema", mock.Anything).Return(&models.Schema{
		Databases: []models.DatabaseSchema{{Tables: []*models.TableSchema{{Name: "users", Engine: "MockEngine"}}}},
	}, nil)
	// Expect query translation
	mockP.On("TranslateQuery", mock.Anything).Return("SELECT * FROM users WHERE id = ?", nil)
	// Expect query execution
	mockP.On("ExecuteQuery", mock.Anything, mock.Anything).Return(&proto.ExecuteQueryResponse{}, nil)

	registry := plugin_registry.NewRegistry()
	registry.Register(mockP)

	parser := parsers.NewAntlrParser()
	convSvc := conversion.NewService(parser, registry)
	genSvc := generation.NewService()
	execSvc := execution.NewService(registry)
	valSvc := validation.NewService()

	mgr := workflow.NewManager(convSvc, genSvc, execSvc, valSvc)

	// 3. Configure Pipeline
	cfg := workflow.WorkflowConfig{
		InputTracePath:      tracePath,
		InputSchemaPath:     schemaPath,
		OutputDir:           filepath.Join(tmpDir, "output"),
		TargetPlugin:        "mock_plugin",
		Generation:          generation.GenerateRequest{Count: 10},
		Execution:           execution.ExecutionConfig{TargetQPS: 100, Concurrency: 1},
		BaselineMetricsPath: baselinePath,
	}

	// 4. Run Pipeline
	err = mgr.Run(context.Background(), cfg)
	require.NoError(t, err)

	// 5. Verify Outputs
	outputDir := cfg.OutputDir
	assert.FileExists(t, filepath.Join(outputDir, "converted", "traces.jsonl"))
	assert.FileExists(t, filepath.Join(outputDir, "converted", "schema.sql"))
	assert.FileExists(t, filepath.Join(outputDir, "workload", "benchmark.jsonl"))
	assert.FileExists(t, filepath.Join(outputDir, "results", "metrics.json"))
	assert.FileExists(t, filepath.Join(outputDir, "report.json"))

	// Verify Schema Conversion result contains MockEngine
	schemaBytes, _ := os.ReadFile(filepath.Join(outputDir, "converted", "schema.sql"))
	assert.Contains(t, string(schemaBytes), "ENGINE = MockEngine")
}
