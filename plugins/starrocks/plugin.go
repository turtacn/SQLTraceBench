package starrocks

import "github.com/turtacn/SQLTraceBench/pkg/plugins"

// PluginName exported
const PluginName = "starrocks"

type plugin struct{}

func New() plugins.DatabasePlugin { return &plugin{} }

func (p *plugin) GetName() string                          { return PluginName }
func (p *plugin) GetVersion() string                       { return "1.0-mvp" }
func (p *plugin) ValidateConnection() error                { return nil }
func (p *plugin) SchemaConverter() plugins.SchemaConverter { return nil }
func (p *plugin) QueryTranslator() plugins.QueryTranslator { return nil }
func (p *plugin) BenchmarkRunner() plugins.BenchmarkRunner { return nil }
