package postgres

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/plugins"
)

const PluginName = "postgres"

type plugin struct{}

func New() plugins.Plugin { return &plugin{} }

func (p *plugin) Name() string    { return PluginName }
func (p *plugin) Version() string { return "1.0-mvp" }
func (p *plugin) TranslateQuery(sql string) (string, error) {
	// Placeholder implementation
	return sql, nil
}

func (p *plugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	// Placeholder implementation for Postgres
	return schema, nil
}
