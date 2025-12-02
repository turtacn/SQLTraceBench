package plugins

import (
	"context"
	"fmt"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// Plugin is the interface that all plugins must implement.
type Plugin interface {
	Name() string
	Version() string
	TranslateQuery(sql string) (string, error)
	ConvertSchema(schema *models.Schema) (*models.Schema, error)
	ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error)
}

// Registry holds a collection of all registered plugins.
type Registry struct {
	plugins map[string]Plugin
}

// NewRegistry creates a new plugin registry.
func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]Plugin)}
}

// Register adds a new plugin to the registry.
func (r *Registry) Register(p Plugin) {
	r.plugins[p.Name()] = p
}

// Get retrieves a plugin from the registry by name.
func (r *Registry) Get(name string) (Plugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}

// GlobalRegistry is the global plugin registry.
var GlobalRegistry = NewRegistry()

// GetPlugin retrieves a plugin from the global registry by name.
func GetPlugin(name string) (Plugin, error) {
	if p, ok := GlobalRegistry.Get(name); ok {
		return p, nil
	}
	return nil, fmt.Errorf("plugin not found: %s", name)
}