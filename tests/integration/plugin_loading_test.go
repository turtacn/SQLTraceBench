package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func TestPluginLoading(t *testing.T) {
	// 1. Find Repo Root
	rootDir, err := findRepoRoot()
	require.NoError(t, err)

	// 2. Build the clickhouse plugin
	pluginDir := filepath.Join(rootDir, "bin", "plugins")
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(filepath.Join(rootDir, "bin")) // Clean up after test

	pluginBinary := filepath.Join(pluginDir, "clickhouse-plugin")

	// Build command - Run from root
	cmd := exec.Command("go", "build", "-o", pluginBinary, "./cmd/plugins/clickhouse")
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	require.NoError(t, err, "Failed to build plugin")

	// 3. Initialize registry and load plugin
	registry := plugin_registry.NewRegistry()
	defer registry.Close()

	err = registry.LoadPluginsFromDir(pluginDir)
	require.NoError(t, err, "Failed to load plugins")

	// 4. Verify plugin availability
	p, ok := registry.Get("clickhouse")
	assert.True(t, ok, "Plugin 'clickhouse' not found in registry")

	if ok {
		assert.NotNil(t, p)
		// 5. Call method
		sql := "SELECT 1"
		translated, err := p.TranslateQuery(sql)
		assert.NoError(t, err)
		// The current clickhouse implementation is a stub returning original SQL
		assert.Equal(t, sql, translated)
	}
}
