package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	// Simulate some work
	return &proto.ExecuteQueryResponse{DurationMicros: 1000}, nil
}

func TestCliPipeline(t *testing.T) {
	// Setup: Create a temporary directory for the test files
	tmpDir, err := ioutil.TempDir("", "cli_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Setup: Register the mock plugin for in-process use
	plugin_registry.GlobalRegistry.Register(&mockPlugin{})
	defer func() {
		// Clean up global state
		plugin_registry.GlobalRegistry = plugin_registry.NewRegistry()
	}()

	// Setup: Create a dummy trace file
	tracePath := filepath.Join(tmpDir, "traces.jsonl")
	traceContent := `[{"timestamp": "2025-01-01T12:00:00Z", "query": "SELECT * FROM users WHERE id = 1", "latency": 120000000}]`
	err = ioutil.WriteFile(tracePath, []byte(traceContent), 0644)
	require.NoError(t, err)

	// Setup: Create a dummy config file to prevent default loading errors
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = ioutil.WriteFile(configPath, []byte(""), 0644)
	require.NoError(t, err)

	// Test the convert command
	tplPath := filepath.Join(tmpDir, "templates.json")
	rootCmd.SetArgs([]string{"--config", configPath, "convert", "--trace", tracePath, "--out", tplPath})
	err = rootCmd.Execute()
	require.NoError(t, err, "convert command failed")
	assert.FileExists(t, tplPath)

	// Test the generate command
	workloadPath := filepath.Join(tmpDir, "workload.json")
	rootCmd.SetArgs([]string{"--config", configPath, "generate", "--source-traces", tracePath, "--out", workloadPath})
	err = rootCmd.Execute()
	require.NoError(t, err, "generate command failed")
	assert.FileExists(t, workloadPath)

	// Test the run command (for base and candidate)
	baseMetricsPath := filepath.Join(tmpDir, "base_metrics.json")
	candMetricsPath := filepath.Join(tmpDir, "cand_metrics.json")

	rootCmd.SetArgs([]string{"--config", configPath, "run", "--workload", workloadPath, "--out", baseMetricsPath, "--db", "mock"})
	err = rootCmd.Execute()
	require.NoError(t, err, "first run command failed")
	assert.FileExists(t, baseMetricsPath)

	rootCmd.SetArgs([]string{"--config", configPath, "run", "--workload", workloadPath, "--out", candMetricsPath, "--db", "mock"})
	err = rootCmd.Execute()
	require.NoError(t, err, "second run command failed")
	assert.FileExists(t, candMetricsPath)

	// Test the validate command
	reportDir := filepath.Join(tmpDir, "report_dir")
	os.MkdirAll(reportDir, 0755)
	rootCmd.SetArgs([]string{"--config", configPath, "validate", "--base", baseMetricsPath, "--candidate", candMetricsPath, "--out", reportDir})
	err = rootCmd.Execute()
	require.NoError(t, err, "validate command failed")
	assert.FileExists(t, filepath.Join(reportDir, "validation_report.html"))

	htmlContent, err := ioutil.ReadFile(filepath.Join(reportDir, "validation_report.html"))
	require.NoError(t, err)
	assert.Contains(t, string(htmlContent), "Validation Report")
}
