package grpc_impl

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// --- Mock Plugin ---

type mockPlugin struct{}

func (m *mockPlugin) Name() string    { return "mock_plugin" }
func (m *mockPlugin) Version() string { return "0.0.1" }
func (m *mockPlugin) TranslateQuery(sql string) (string, error) {
	return "TRANSLATED: " + sql, nil
}

// ConvertSchema implements the SchemaConverter interface required by server.go logic
func (m *mockPlugin) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	// For testing, we modify the schema to prove it was processed as an object
	if len(schema.Databases) > 0 {
		schema.Databases[0].Name = "CONVERTED_" + schema.Databases[0].Name
	}
	return schema, nil
}

func (m *mockPlugin) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	return &proto.ExecuteQueryResponse{}, nil
}

// --- Test Mapper ---

func TestMapper(t *testing.T) {
	// Create a complex schema
	original := &models.Schema{
		Databases: []models.DatabaseSchema{
			{
				Name: "db1",
				Tables: []*models.TableSchema{
					{
						Name: "table1",
						Columns: []*models.ColumnSchema{
							{Name: "col1", DataType: "INT", IsPrimaryKey: true},
							{Name: "col2", DataType: "VARCHAR", IsNullable: true},
						},
						PK: []string{"col1"},
						Indexes: map[string]*models.IndexSchema{
							"idx1": {Name: "idx1", Columns: []string{"col2"}, IsUnique: false},
						},
					},
				},
			},
		},
	}

	// Test ToProtoSchema
	jsonStr, err := ToProtoSchema(original)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	// Test FromProtoSchema
	restored, err := FromProtoSchema(jsonStr)
	assert.NoError(t, err)
	assert.Equal(t, original, restored)
}

func TestTraceMapper(t *testing.T) {
	trace := &models.SQLTrace{Query: "SELECT * FROM users"}

	// Trace -> Query
	query := TraceToQuery(trace)
	assert.Equal(t, "SELECT * FROM users", query)

	// Query -> Trace
	newTrace := QueryToTrace(query)
	assert.Equal(t, trace.Query, newTrace.Query)
}

// --- Test GRPC Flow ---

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPCFlow(t *testing.T) {
	// 1. Setup Server
	s := grpc.NewServer()
	mockImpl := &mockPlugin{}

	// Register the GRPCServer with our mock implementation
	proto.RegisterSQLTraceBenchPluginServer(s, &GRPCServer{Impl: mockImpl})

	go func() {
		if err := s.Serve(lis); err != nil {
			// server might be closed, ignore
		}
	}()
	defer s.Stop()

	// 2. Setup Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	// Create GRPCClient manually
	client := &GRPCClient{
		client: proto.NewSQLTraceBenchPluginClient(conn),
	}

	// 3. Test TranslateQuery
	sql := "SELECT * FROM t"
	translated, err := client.TranslateQuery(sql)
	assert.NoError(t, err)
	assert.Equal(t, "TRANSLATED: SELECT * FROM t", translated)

	// 4. Test ConvertSchema with domain objects
	schemaObj := &models.Schema{Databases: []models.DatabaseSchema{{Name: "test_db"}}}

	convertedSchema, err := client.ConvertSchema(schemaObj)
	assert.NoError(t, err)
	assert.NotNil(t, convertedSchema)
	assert.Equal(t, "CONVERTED_test_db", convertedSchema.Databases[0].Name)

	// 5. Test GetName
	name := client.GetName()
	assert.Equal(t, "mock_plugin", name)

	// Test Name() alias
	assert.Equal(t, "mock_plugin", client.Name())

	// Test Version()
	assert.Equal(t, "1.0.0", client.Version())
}
