// Package storage provides in-memory storage implementations for the application's domain models.
package storage

import (
	"sync"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// MemoryStorage provides an in-memory implementation of a key-value store for SQL templates.
// It uses a mutex to ensure thread-safe access to the underlying map.
type MemoryStorage struct {
	tplMu sync.RWMutex
	tpls  map[string]*models.SQLTemplate
}

// NewMemoryStorage creates a new MemoryStorage instance.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{tpls: make(map[string]*models.SQLTemplate)}
}

// SaveTemplate saves a SQL template to the in-memory store.
// It uses the template's GroupKey as the key.
func (m *MemoryStorage) SaveTemplate(t *models.SQLTemplate) {
	m.tplMu.Lock()
	defer m.tplMu.Unlock()
	m.tpls[t.GroupKey] = t
}

// GetTemplate retrieves a SQL template from the in-memory store by its ID (GroupKey).
func (m *MemoryStorage) GetTemplate(id string) (*models.SQLTemplate, bool) {
	m.tplMu.RLock()
	defer m.tplMu.RUnlock()
	t, ok := m.tpls[id]
	return t, ok
}