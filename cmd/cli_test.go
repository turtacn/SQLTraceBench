package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestCliPipeline(t *testing.T) {
	// Setup: Create a temporary directory for the test files
	tmpDir, err := ioutil.TempDir("", "cli_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Setup: Create a dummy trace file
	tracePath := filepath.Join(tmpDir, "traces.jsonl")
	traceContent := `
{"query": "select * from users where id = 1"}
{"query": "select * from users where id = 2"}
{"query": "select * from users where id = 1"}
`
	ioutil.WriteFile(tracePath, []byte(traceContent), 0644)

	// Test the convert command
	tplPath := filepath.Join(tmpDir, "templates.json")
	rootCmd.SetArgs([]string{"convert", "--trace", tracePath, "--out", tplPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, tplPath)

	// Test the generate command
	workloadPath := filepath.Join(tmpDir, "workload.json")
	rootCmd.SetArgs([]string{"generate", "--templates", tplPath, "--out", workloadPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, workloadPath)

	// Test the run command (for base and candidate)
	baseMetricsPath := filepath.Join(tmpDir, "base_metrics.json")
	candMetricsPath := filepath.Join(tmpDir, "cand_metrics.json")
	rootCmd.SetArgs([]string{"run", "--workload", workloadPath, "--out", baseMetricsPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, baseMetricsPath)
	rootCmd.SetArgs([]string{"run", "--workload", workloadPath, "--out", candMetricsPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, candMetricsPath)

	// Test the validate command
	reportPath := filepath.Join(tmpDir, "report.json")
	rootCmd.SetArgs([]string{"validate", "--base", baseMetricsPath, "--candidate", candMetricsPath, "--out", reportPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, reportPath)

	// Verify the report
	var report models.Report
	reportData, err := ioutil.ReadFile(reportPath)
	require.NoError(t, err)
	err = json.Unmarshal(reportData, &report)
	require.NoError(t, err)
	assert.True(t, report.Result.Pass, "validation should pass for similar runs")
}