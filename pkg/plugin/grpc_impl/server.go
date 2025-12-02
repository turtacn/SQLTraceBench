package grpc_impl

import (
	"context"

	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"github.com/turtacn/SQLTraceBench/plugins"
)

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
	// Deserialize input string to domain object
	domainSchema, err := FromProtoSchema(req.Schema)
	if err != nil {
		return &proto.ConvertSchemaResponse{Error: "invalid schema format: " + err.Error()}, nil
	}

	// Call implementation directly as plugins.Plugin interface now includes ConvertSchema
	resSchema, err := s.Impl.ConvertSchema(domainSchema)
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

func (s *GRPCServer) ExecuteQuery(ctx context.Context, req *proto.ExecuteQueryRequest) (*proto.ExecuteQueryResponse, error) {
	return s.Impl.ExecuteQuery(ctx, req)
}
