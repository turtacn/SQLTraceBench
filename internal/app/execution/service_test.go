package execution

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"github.com/turtacn/SQLTraceBench/plugins"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

type MockPlugin struct {
	plugins.Plugin
}

func (p *MockPlugin) Name() string {
	return "mock"
}

func (p *MockPlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	return &proto.ExecuteQueryResponse{}, nil
}

func TestDefaultService_RunBenchmark(t *testing.T) {
	// Create a mock plugin registry.
	registry := plugin_registry.NewRegistry()
	registry.Register(&MockPlugin{})

	// Create the service.
	service := NewService(registry)

	// Create a benchmark workload.
	workload := &models.BenchmarkWorkload{
		Queries: []models.QueryWithArgs{
			{Query: "SELECT 1", Args: []interface{}{}},
		},
	}

	// Create an execution config.
	config := ExecutionConfig{
		TargetDB:    "mock",
		TargetQPS:   10,
		Concurrency: 1,
	}

	// Run the benchmark.
	result, err := service.RunBenchmark(context.Background(), workload, config)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, result)
}