package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Setup: Create a temporary config file
	tmpDir, err := ioutil.TempDir("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
log:
  level: debug
  format: json
database:
  dsn: "user:pass@tcp(localhost:3306)/testdb"
benchmark:
  qps: 200
`
	ioutil.WriteFile(configPath, []byte(configContent), 0644)

	// Test loading from the file
	cfg, err := Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "user:pass@tcp(localhost:3306)/testdb", cfg.Database.DSN)
	assert.Equal(t, 200, cfg.Benchmark.QPS)
	assert.Equal(t, 10, cfg.Benchmark.Concurrency) // Default value

	// Test environment variable override
	os.Setenv("SQLTRACEBENCH_BENCHMARK_CONCURRENCY", "50")
	os.Setenv("SQLTRACEBENCH_DATABASE_DRIVER", "postgres")
	cfg, err = Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, 50, cfg.Benchmark.Concurrency)
	assert.Equal(t, "postgres", cfg.Database.Driver)
}