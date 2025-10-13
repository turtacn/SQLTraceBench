package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	// Create a new plugin.
	p := New()

	// Test the Name and Version methods.
	assert.Equal(t, "clickhouse", p.Name())
	assert.Equal(t, "1.0-mvp", p.Version())

	// Test the TranslateQuery method.
	translated, err := p.TranslateQuery("SELECT 1")
	assert.NoError(t, err)
	assert.Equal(t, "SELECT 1", translated)
}