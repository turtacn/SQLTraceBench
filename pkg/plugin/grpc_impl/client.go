package grpc_impl

import (
	"context"
	"errors"

	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// GRPCClient is an implementation of plugins.Plugin that talks over RPC.
type GRPCClient struct {
	client proto.SQLTraceBenchPluginClient
}

func (c *GRPCClient) Name() string {
	resp, err := c.client.GetName(context.Background(), &proto.Empty{})
	if err != nil {
		return "unknown"
	}
	return resp.Name
}

func (c *GRPCClient) Version() string {
	// Version is not in proto yet.
	return "1.0.0"
}

func (c *GRPCClient) TranslateQuery(sql string) (string, error) {
	resp, err := c.client.TranslateQuery(context.Background(), &proto.TranslateQueryRequest{Sql: sql})
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}
	return resp.TranslatedSql, nil
}
