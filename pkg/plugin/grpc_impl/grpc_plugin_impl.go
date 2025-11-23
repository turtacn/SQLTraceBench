package grpc_impl

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
	"github.com/turtacn/SQLTraceBench/plugins"
	"google.golang.org/grpc"
)

// GRPCPluginImpl implements plugin.GRPCPlugin interface.
// It's the bridge between Hashicorp plugin system and our gRPC implementation.
type GRPCPluginImpl struct {
	plugin.Plugin
	// Impl is the actual implementation of the business logic (server side).
	Impl plugins.Plugin
}

func (p *GRPCPluginImpl) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterSQLTraceBenchPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *GRPCPluginImpl) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewSQLTraceBenchPluginClient(c)}, nil
}
