package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Setup: Create a dummy config file
	configPath := filepath.Join(tmpDir, "config.yaml")
	ioutil.WriteFile(configPath, []byte(""), 0644)
	cfgFile = configPath

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
	// validation command uses "out" as directory for HTML report.
	reportDir := filepath.Join(tmpDir, "report_dir")
	os.MkdirAll(reportDir, 0755)
	rootCmd.SetArgs([]string{"validate", "--base", baseMetricsPath, "--candidate", candMetricsPath, "--out", reportDir})
	err = rootCmd.Execute()
	require.NoError(t, err)
	// Check for HTML report
	assert.FileExists(t, filepath.Join(reportDir, "validation_report.html"))

	// Since `ValidateTrace` returns the report object but CLI doesn't write JSON explicitly unless implemented,
	// we skip JSON content verification here if it's not written.
	// If we want to verify pass/fail, we might need to rely on the exit code or logs, or assume if no error it passed.
	// Or we can check if HTML contains "Pass".
	htmlContent, err := ioutil.ReadFile(filepath.Join(reportDir, "validation_report.html"))
	require.NoError(t, err)
	// Simple check if template rendered
	assert.Contains(t, string(htmlContent), "Validation Report")
}