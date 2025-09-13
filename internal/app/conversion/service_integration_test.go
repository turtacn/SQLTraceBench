package conversion

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// TestTemplate is a simplified struct for comparison, ignoring dynamic fields like ID and CreatedAt
type TestTemplate struct {
	TemplateQuery string
	Frequency     int64
}

func TestConvertFromFile_Integration(t *testing.T) {
	// Setup: Create the service with its dependencies
	templateSvc := services.NewTemplateService(nil) // No repo needed for this test
	schemaSvc := services.NewSchemaService()
	service := NewService(templateSvc, schemaSvc)

	// Setup: Define file paths
	testdataDir := "testdata"
	tracePath := filepath.Join(testdataDir, "source_traces.jsonl")
	expectedJSONPath := filepath.Join(testdataDir, "expected_templates.json")

	// Setup: Create a temporary file for the output
	tmpfile, err := ioutil.TempFile("", "test_output_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfilePath := tmpfile.Name()
	tmpfile.Close() // Close the file so the service can create/write to it

	// Execute the service method
	err = service.ConvertFromFile(context.Background(), tracePath, tmpfilePath)
	require.NoError(t, err)

	// Verify: Read and unmarshal the actual output
	actualData, err := ioutil.ReadFile(tmpfilePath)
	require.NoError(t, err)

	var actualResult struct {
		Templates []models.SQLTemplate `json:"templates"`
	}
	err = json.Unmarshal(actualData, &actualResult)
	require.NoError(t, err)

	// Verify: Read and unmarshal the expected output
	expectedData, err := ioutil.ReadFile(expectedJSONPath)
	require.NoError(t, err)

	var expectedResult struct {
		Templates []TestTemplate `json:"templates"`
	}
	err = json.Unmarshal(expectedData, &expectedResult)
	require.NoError(t, err)

	// Verify: Compare the results in an order-insensitive way
	// Convert actual templates to the simplified TestTemplate struct for comparison
	actualTemplatesForTest := make(map[string]TestTemplate)
	for _, tmpl := range actualResult.Templates {
		actualTemplatesForTest[tmpl.TemplateQuery] = TestTemplate{
			TemplateQuery: tmpl.TemplateQuery,
			Frequency:     tmpl.Frequency,
		}
	}

	expectedTemplatesForTest := make(map[string]TestTemplate)
	for _, tmpl := range expectedResult.Templates {
		expectedTemplatesForTest[tmpl.TemplateQuery] = tmpl
	}

	assert.Equal(t, expectedTemplatesForTest, actualTemplatesForTest, "The generated templates should match the expected templates")
}
