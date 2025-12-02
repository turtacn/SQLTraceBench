package clickhouse

import (
	"context"
	"database/sql"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// ClickHousePlugin implements the DatabasePlugin interface.
type ClickHousePlugin struct {
	converter  SchemaConverter
	translator QueryTranslator
	executor   *BenchmarkExecutor
}

// New creates a new ClickHousePlugin instance.
func New() *ClickHousePlugin {
	// For now, we'll create a new connection each time.
	// In a real application, you'd use a connection pool.
	conn, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		// In a real application, you'd handle this error more gracefully.
		// For now, we'll panic.
		panic(err)
	}
	return &ClickHousePlugin{
		converter:  NewSchemaConverter(),
		translator: NewQueryTranslator(),
		executor:   NewBenchmarkExecutor(conn),
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

// GetBenchmarkExecutor returns the benchmark executor.
func (p *ClickHousePlugin) GetBenchmarkExecutor() (*BenchmarkExecutor, error) {
	// For now, we'll create a new connection each time.
	// In a real application, you'd use a connection pool.
	conn, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		return nil, err
	}
	p.executor = NewBenchmarkExecutor(conn)
	return p.executor, nil
}

func (p *ClickHousePlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	return p.executor.ExecuteQuery(ctx, req)
}
