package storage

import (
	"sync"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type MemoryStorage struct {
	tplMu sync.RWMutex
	tpls  map[string]*models.SQLTemplate
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{tpls: make(map[string]*models.SQLTemplate)}
}

func (m *MemoryStorage) SaveTemplate(t *models.SQLTemplate) {
	m.tplMu.Lock()
	defer m.tplMu.Unlock()
	m.tpls[t.TemplateID] = t
}

func (m *MemoryStorage) GetTemplate(id string) (*models.SQLTemplate, bool) {
	m.tplMu.RLock()
	defer m.tplMu.RUnlock()
	t, ok := m.tpls[id]
	return t, ok
}

//Personal.AI order the ending
