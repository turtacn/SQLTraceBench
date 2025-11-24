package clickhouse

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ClickHousePlugin implements the DatabasePlugin interface.
type ClickHousePlugin struct {
	converter  SchemaConverter
	translator QueryTranslator
}

// New creates a new ClickHousePlugin instance.
func New() *ClickHousePlugin {
	return &ClickHousePlugin{
		converter:  NewSchemaConverter(),
		translator: NewQueryTranslator(),
	}
}

// Name returns the name of the plugin.
func (p *ClickHousePlugin) Name() string {
	return "clickhouse" // Should match GetName() requirement AC-2
}

// GetName returns the name of the plugin (alias for Name to match interface if needed, or interface uses GetName)
// Interface definition: GetName() string
func (p *ClickHousePlugin) GetName() string {
	return "clickhouse"
}

// Version returns the version of the plugin.
func (p *ClickHousePlugin) Version() string {
	return "1.0.0"
}

// TranslateQuery translates a SQL query.
func (p *ClickHousePlugin) TranslateQuery(sql string) (string, error) {
	return p.translator.TranslateQuery(sql)
}

// ConvertSchema converts a database schema.
func (p *ClickHousePlugin) ConvertSchema(src *models.Schema) (*models.Schema, error) {
	return p.converter.ConvertSchema(src)
}
