package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

// mockPlugin is a mock implementation of the plugins.Plugin interface for in-process testing.
type mockPlugin struct{}

func (m *mockPlugin) Name() string {
	return "mock"
}

func (m *mockPlugin) Version() string {
	return "0.1.0"
}

func (m *mockPlugin) TranslateQuery(query string) (string, error) {
	return query, nil
}

func (m *mockPlugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	return schema, nil
}

func (m *mockPlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	time.Sleep(10 * time.Millisecond)
	return &proto.ExecuteQueryResponse{DurationMicros: 10000}, nil
}

func TestExecutionFlow(t *testing.T) {
	// Create a mock plugin registry and register our in-process mock.
	registry := plugin_registry.NewRegistry()
	registry.Register(&mockPlugin{})

	// Create the service.
	service := execution.NewService(registry)

	// Create a benchmark workload.
	workload := &models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}

	// Create an execution config.
	config := execution.ExecutionConfig{
		TargetDB:    "mock",
		TargetQPS:   10,
		Concurrency: 1,
	}

	// Run the benchmark.
	result, err := service.RunBenchmark(context.Background(), workload, config)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
