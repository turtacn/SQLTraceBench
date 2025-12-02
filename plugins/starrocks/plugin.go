package starrocks

import (
	"context"
	"database/sql"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// PluginName exported
const PluginName = "starrocks"
const PluginVersion = "1.0.0"

// StarRocksPlugin implements the DatabasePlugin interface.
type StarRocksPlugin struct {
	converter  *StarRocksConverter
	translator *StarRocksTranslator
	executor   *BenchmarkExecutor
}

// New creates a new instance of StarRocksPlugin.
func New() *StarRocksPlugin {
	return &StarRocksPlugin{
		converter:  NewSchemaConverter(),
		translator: &StarRocksTranslator{},
	}
}

// Name returns the name of the plugin.
func (p *StarRocksPlugin) Name() string {
	return PluginName
}

// GetName alias for interface compatibility
func (p *StarRocksPlugin) GetName() string {
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

// GetBenchmarkExecutor returns the benchmark executor.
func (p *StarRocksPlugin) GetBenchmarkExecutor() (*BenchmarkExecutor, error) {
	// For now, we'll create a new connection each time.
	// In a real application, you'd use a connection pool.
	conn, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:9030)/dbname")
	if err != nil {
		return nil, err
	}
	p.executor = NewBenchmarkExecutor(conn)
	return p.executor, nil
}

func (p *StarRocksPlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	return p.executor.ExecuteQuery(ctx, req)
}
