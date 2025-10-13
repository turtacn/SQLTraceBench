package plugin_registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/plugins"
)

func TestPluginRegistry(t *testing.T) {
	// The init function in this package should have registered the plugins.
	// We can test this by trying to get them from the global registry.

	// Test that the clickhouse plugin is registered.
	_, err := plugins.GetPlugin("clickhouse")
	assert.NoError(t, err)

	// Test that the starrocks plugin is registered.
	_, err = plugins.GetPlugin("starrocks")
	assert.NoError(t, err)

	// Test that the postgres plugin is registered.
	_, err = plugins.GetPlugin("postgres")
	assert.NoError(t, err)
}