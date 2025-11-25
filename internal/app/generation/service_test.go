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
	// Check if queries are generated with arguments
	for _, q := range workload.Queries {
		assert.NotEmpty(t, q.Query)
		assert.NotNil(t, q.Args)
	}
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

	// Prepare a request with traces that might not yield templates (e.g., malformed).
	// For this test, let's use valid traces but imagine a scenario where template extraction fails.
	// We can't easily mock the template service here, so we test the path where templates might be empty.
	// A trace with no query will result in no templates.
	req := GenerateRequest{
		Count: 10,
		SourceTraces: []models.SQLTrace{
			{Query: ""},
		},
	}

	// Generate the workload.
	_, err := service.GenerateWorkload(context.Background(), req)

	// Assert the error.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no templates could be extracted")
}
