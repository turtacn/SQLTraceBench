package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestMemoryStorage(t *testing.T) {
	// Create a new memory store.
	store := NewMemoryStorage()

	// Test saving and getting a template.
	tpl1 := &models.SQLTemplate{GroupKey: "key1", RawSQL: "SELECT 1"}
	store.SaveTemplate(tpl1)
	retrieved, ok := store.GetTemplate("key1")
	assert.True(t, ok)
	assert.Equal(t, tpl1, retrieved)

	// Test getting a non-existent template.
	_, ok = store.GetTemplate("key2")
	assert.False(t, ok)

	// Test saving multiple templates.
	tpl2 := &models.SQLTemplate{GroupKey: "key2", RawSQL: "SELECT 2"}
	store.SaveTemplate(tpl2)
	retrieved, ok = store.GetTemplate("key1")
	assert.True(t, ok)
	assert.Equal(t, tpl1, retrieved)
	retrieved, ok = store.GetTemplate("key2")
	assert.True(t, ok)
	assert.Equal(t, tpl2, retrieved)
}