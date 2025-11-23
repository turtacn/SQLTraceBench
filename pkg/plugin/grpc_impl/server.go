package grpc_impl

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"github.com/turtacn/SQLTraceBench/plugins"
)

// SchemaConverter is an interface that plugins can optionally implement if they support schema conversion.
type SchemaConverter interface {
	ConvertSchema(schema *models.Schema) (*models.Schema, error)
}

// GRPCServer implements the proto.SQLTraceBenchPluginServer interface.
// It receives gRPC requests and forwards them to the actual implementation (plugins.Plugin).
type GRPCServer struct {
	proto.UnimplementedSQLTraceBenchPluginServer
	Impl plugins.Plugin
}

func (s *GRPCServer) GetName(ctx context.Context, req *proto.Empty) (*proto.NameResponse, error) {
	return &proto.NameResponse{Name: s.Impl.Name()}, nil
}

func (s *GRPCServer) TranslateQuery(ctx context.Context, req *proto.TranslateQueryRequest) (*proto.TranslateQueryResponse, error) {
	// Directly call TranslateQuery as per current plugins.Plugin interface
	res, err := s.Impl.TranslateQuery(req.Sql)
	if err != nil {
		return &proto.TranslateQueryResponse{Error: err.Error()}, nil
	}
	return &proto.TranslateQueryResponse{TranslatedSql: res}, nil
}

func (s *GRPCServer) ConvertSchema(ctx context.Context, req *proto.ConvertSchemaRequest) (*proto.ConvertSchemaResponse, error) {
	// Check if the plugin implementation supports SchemaConverter
	if converter, ok := s.Impl.(SchemaConverter); ok {
		// Deserialize input string to domain object
		domainSchema, err := FromProtoSchema(req.Schema)
		if err != nil {
			return &proto.ConvertSchemaResponse{Error: "invalid schema format: " + err.Error()}, nil
		}

		// Call implementation
		resSchema, err := converter.ConvertSchema(domainSchema)
		if err != nil {
			return &proto.ConvertSchemaResponse{Error: err.Error()}, nil
		}

		// Serialize result domain object back to string
		resStr, err := ToProtoSchema(resSchema)
		if err != nil {
			return &proto.ConvertSchemaResponse{Error: "failed to serialize response: " + err.Error()}, nil
		}

		return &proto.ConvertSchemaResponse{ConvertedSchema: resStr}, nil
	}

	return &proto.ConvertSchemaResponse{Error: "schema conversion not supported by this plugin"}, nil
}
