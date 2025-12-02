package generation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestDefaultService_GenerateWorkload_Success(t *testing.T) {
	// Create the service.
	service := NewService()

	// Prepare a request with source traces.
	req := GenerateRequest{
		Count: 10,
		SourceTraces: []models.SQLTrace{
			{
				Query:      "SELECT * FROM users WHERE id = ?",
				Parameters: map[string]interface{}{"id": 1},
			},
			{
				Query:      "SELECT * FROM products WHERE sku = ?",
				Parameters: map[string]interface{}{"sku": "ABC"},
			},
		},
	}

	// Generate the workload.
	workload, err := service.GenerateWorkload(context.Background(), req)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, workload)
	assert.Len(t, workload.Queries, 10)
}

func TestDefaultService_GenerateWorkload_NoTraces(t *testing.T) {
	// Create the service.
	service := NewService()

	// Prepare a request with no source traces.
	req := GenerateRequest{
		Count:        10,
		SourceTraces: []models.SQLTrace{},
	}

	// Generate the workload.
	_, err := service.GenerateWorkload(context.Background(), req)

	// Assert the error.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "generation requires source traces")
}

func TestDefaultService_GenerateWorkload_NoTemplatesExtracted(t *testing.T) {
	// Create the service.
	service := NewService()

	// Prepare a request with traces that might not yield parameterized templates.
	// The service should still treat them as templates and generate a workload.
	req := GenerateRequest{
		Count: 10,
		SourceTraces: []models.SQLTrace{
			{Query: "SELECT 1"},
		},
	}

	// Generate the workload.
	workload, err := service.GenerateWorkload(context.Background(), req)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, workload)
	assert.Len(t, workload.Queries, 10) // Expect 10 queries, even if no parameters are extracted
}
