package plugin_registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	pkg_plugin "github.com/turtacn/SQLTraceBench/pkg/plugin"
	"github.com/turtacn/SQLTraceBench/plugins"
)

// Registry holds a collection of all registered plugins.
// For gRPC plugins, it manages the go-plugin Clients.
type Registry struct {
	plugins map[string]plugins.Plugin
	clients []*plugin.Client
}

// GlobalRegistry is the global plugin registry.
var GlobalRegistry *Registry

func init() {
	GlobalRegistry = NewRegistry()
}

// NewRegistry creates a new plugin registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]plugins.Plugin),
		clients: make([]*plugin.Client, 0),
	}
}

// LoadPluginsFromDir scans the directory for plugin binaries and loads them.
func (r *Registry) LoadPluginsFromDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		// If directory doesn't exist, we might just warn or return error depending on strictness.
		// For now, if it doesn't exist, it's just empty.
		if os.IsNotExist(err) {
			logrus.Warnf("Plugin directory %s does not exist", dir)
			return nil
		}
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		// Simple convention: any file in bin/plugins (or --plugin-dir) is a potential plugin.
		// We might want to filter by prefix or executable bit.
		// For Phase 1, let's assume anything executable or with a specific suffix.
		// Let's assume the filename is the plugin name for now, or we just try to load it.
		// The prompt says: "LoadPlugin 函数接收文件路径"

		fullPath := filepath.Join(dir, f.Name())
		// Skip non-executable files (basic check)
		info, err := f.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0111 == 0 {
			continue
		}

		if err := r.LoadPlugin(fullPath); err != nil {
			logrus.Warnf("Failed to load plugin %s: %v", f.Name(), err)
		} else {
			logrus.Infof("Loaded plugin from %s", f.Name())
		}
	}
	return nil
}

// LoadPlugin loads a single plugin from the given path.
func (r *Registry) LoadPlugin(path string) error {
	// Create an hclog.Logger that writes to stderr or compatible place
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: pkg_plugin.HandshakeConfig,
		Plugins:         pkg_plugin.PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:          logger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create rpc client: %w", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("database_plugin")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin: %w", err)
	}

	dbPlugin, ok := raw.(plugins.Plugin)
	if !ok {
		client.Kill()
		return fmt.Errorf("plugin does not implement plugins.Plugin interface")
	}

	r.plugins[dbPlugin.Name()] = dbPlugin
	r.clients = append(r.clients, client)

	// Also register to the legacy global registry for compatibility
	plugins.GlobalRegistry.Register(dbPlugin)
	return nil
}

// Get retrieves a plugin from the registry by name.
func (r *Registry) Get(name string) (plugins.Plugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}

// Close kills all plugin clients.
func (r *Registry) Close() {
	for _, client := range r.clients {
		client.Kill()
	}
}

// GetPlugin retrieves a plugin from the global registry by name.
func GetPlugin(name string) (plugins.Plugin, error) {
	if p, ok := GlobalRegistry.Get(name); ok {
		return p, nil
	}
	return nil, fmt.Errorf("plugin not found: %s", name)
}

// LoadPlugins loads plugins from the specified directory into the global registry.
func LoadPlugins(dir string) error {
	return GlobalRegistry.LoadPluginsFromDir(dir)
}

// ClosePlugins closes all plugins in the global registry.
func ClosePlugins() {
	GlobalRegistry.Close()
}
