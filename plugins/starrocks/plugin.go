package starrocks

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// PluginName exported
const PluginName = "starrocks"
const PluginVersion = "1.0.0"

// StarRocksPlugin implements the DatabasePlugin interface.
type StarRocksPlugin struct {
	converter  *StarRocksConverter
	translator *StarRocksTranslator
}

// New creates a new instance of StarRocksPlugin.
func New() *StarRocksPlugin {
	return &StarRocksPlugin{
		converter:  &StarRocksConverter{},
		translator: &StarRocksTranslator{},
	}
}

// Name returns the name of the plugin.
func (p *StarRocksPlugin) Name() string {
	return PluginName
}

// Version returns the version of the plugin.
func (p *StarRocksPlugin) Version() string {
	return PluginVersion
}

// TranslateQuery translates a SQL query to StarRocks dialect.
func (p *StarRocksPlugin) TranslateQuery(sql string) (string, error) {
	return p.translator.Translate(sql)
}

// ConvertSchema converts a schema to StarRocks schema.
func (p *StarRocksPlugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	return p.converter.ConvertSchema(schema)
}
