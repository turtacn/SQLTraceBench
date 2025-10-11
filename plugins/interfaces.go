// Package plugins defines the interface for extending the functionality of SQLTraceBench.
package plugins

// Plugin is the interface that all plugins must implement.
// It provides a way to add support for new database dialects.
type Plugin interface {
	// Name returns the name of the plugin (e.g., "clickhouse").
	Name() string
	// Version returns the version of the plugin.
	Version() string
	// TranslateQuery translates a SQL query from the source dialect to the target dialect.
	TranslateQuery(sql string) (string, error)
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