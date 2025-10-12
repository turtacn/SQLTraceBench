package clickhouse

import "github.com/turtacn/SQLTraceBench/plugins"

const PluginName = "clickhouse"

type plugin struct{}

func New() plugins.Plugin { return &plugin{} }

func (p *plugin) Name() string    { return PluginName }
func (p *plugin) Version() string { return "1.0-mvp" }
func (p *plugin) TranslateQuery(sql string) (string, error) {
	// Placeholder implementation
	return sql, nil
}