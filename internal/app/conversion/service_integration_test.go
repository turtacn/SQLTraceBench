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
)

// TestTemplate is a simplified struct for comparison.
type TestTemplate struct {
	GroupKey string `json:"group_key"`
	Weight   int    `json:"weight"`
}

func TestConvertFromFile_Integration(t *testing.T) {
	// Setup: Create the service
	service := NewService()

	// Setup: Define file paths
	testdataDir := "testdata"
	tracePath := filepath.Join(testdataDir, "source_traces.jsonl")
	expectedJSONPath := filepath.Join(testdataDir, "expected_templates.json")

	// Setup: Create test data files
	os.MkdirAll(testdataDir, 0755)
	traceContent := `
{"query": "select * from users where id = :id"}
{"query": "SELECT * FROM users WHERE id = :id "}
{"query": "select * from orders"}
{"query": "select * from users where id = :id"}
`
	ioutil.WriteFile(tracePath, []byte(traceContent), 0644)

	expectedResultTemplates := []TestTemplate{
		{GroupKey: "select * from users where id = :id", Weight: 3},
		{GroupKey: "select * from orders", Weight: 1},
	}
	expectedContent, _ := json.Marshal(map[string][]TestTemplate{"templates": expectedResultTemplates})
	ioutil.WriteFile(expectedJSONPath, expectedContent, 0644)
	defer os.RemoveAll(testdataDir)

	// Execute the service method
	actualTemplates, err := service.ConvertFromFile(context.Background(), tracePath)
	require.NoError(t, err)

	// Verify: Compare the results
	assert.Len(t, actualTemplates, 2)
	foundUsers := false
	foundOrders := false
	for _, tmpl := range actualTemplates {
		if tmpl.GroupKey == "select * from users where id = :id" {
			assert.Equal(t, 3, tmpl.Weight)
			foundUsers = true
		}
		if tmpl.GroupKey == "select * from orders" {
			assert.Equal(t, 1, tmpl.Weight)
			foundOrders = true
		}
	}
	assert.True(t, foundUsers, "users template should be found")
	assert.True(t, foundOrders, "orders template should be found")
}