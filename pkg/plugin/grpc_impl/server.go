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
	// Current plugins.Plugin does not have schema conversion yet.
	// Returning error or empty for now.
	return &proto.ConvertSchemaResponse{Error: "not implemented"}, nil
}
