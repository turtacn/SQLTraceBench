package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type MockPlugin struct {
	name    string
	version string
}

func (p *MockPlugin) Name() string {
	return p.name
}

func (p *MockPlugin) Version() string {
	return p.version
}

func (p *MockPlugin) TranslateQuery(sql string) (string, error) {
	return sql, nil
}

func (p *MockPlugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	return schema, nil
}

func TestRegistry(t *testing.T) {
	// Create a new registry.
	registry := NewRegistry()

	// Test registering and getting a plugin.
	plugin1 := &MockPlugin{name: "plugin1", version: "1.0"}
	registry.Register(plugin1)
	retrieved, ok := registry.Get("plugin1")
	assert.True(t, ok)
	assert.Equal(t, plugin1, retrieved)

	// Test getting a non-existent plugin.
	_, ok = registry.Get("plugin2")
	assert.False(t, ok)
}

func TestGetPlugin(t *testing.T) {
	// Register a plugin to the global registry.
	plugin1 := &MockPlugin{name: "plugin1", version: "1.0"}
	GlobalRegistry.Register(plugin1)

	// Test getting a plugin from the global registry.
	retrieved, err := GetPlugin("plugin1")
	assert.NoError(t, err)
	assert.Equal(t, plugin1, retrieved)

	// Test getting a non-existent plugin from the global registry.
	_, err = GetPlugin("plugin2")
	assert.Error(t, err)
}