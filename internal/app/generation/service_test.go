package generation

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func createTempTemplateFile(t *testing.T, templates []models.SQLTemplate) string {
	t.Helper()
	file, err := os.CreateTemp("", "templates-*.json")
	require.NoError(t, err)
	defer file.Close()

	err = json.NewEncoder(file).Encode(templates)
	require.NoError(t, err)

	return file.Name()
}

func TestDefaultService_GenerateWorkload_Success(t *testing.T) {
	// Create a temporary template file.
	templates := []models.SQLTemplate{
		{
			RawSQL: "SELECT * FROM users WHERE id = :id",
			Weight: 1,
		},
	}
	templatePath := createTempTemplateFile(t, templates)
	defer os.Remove(templatePath)

	// Create the service.
	service := NewService()

	// Generate the workload.
	workload, err := service.GenerateWorkload(context.Background(), templatePath, 10)

	// Assert the results.
	require.NoError(t, err)
	assert.NotNil(t, workload)
	assert.Len(t, workload.Queries, 10)
}

func TestDefaultService_GenerateWorkload_TemplateNotFound(t *testing.T) {
	// Create the service.
	service := NewService()

	// Generate the workload with a non-existent template file.
	_, err := service.GenerateWorkload(context.Background(), "non-existent-file.json", 10)

	// Assert the error.
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestDefaultService_GenerateWorkload_MalformedTemplate(t *testing.T) {
	// Create a malformed template file.
	file, err := os.CreateTemp("", "malformed-*.json")
	require.NoError(t, err)
	defer os.Remove(file.Name())
	_, err = file.WriteString("this is not json")
	require.NoError(t, err)
	file.Close()

	// Create the service.
	service := NewService()

	// Generate the workload.
	_, err = service.GenerateWorkload(context.Background(), file.Name(), 10)

	// Assert the error.
	require.Error(t, err)
}