package grpc_impl

import (
	"context"
	"errors"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

// GRPCClient is an implementation of DatabasePlugin that talks over RPC.
type GRPCClient struct {
	client proto.SQLTraceBenchPluginClient
}

// Name implements plugins.Plugin interface (and DatabasePlugin via GetName wrapper if needed, but the interface says GetName)
// Wait, DatabasePlugin interface has GetName. plugins.Plugin has Name.
// We implement GetName for DatabasePlugin.
func (c *GRPCClient) GetName() string {
	resp, err := c.client.GetName(context.Background(), &proto.Empty{})
	if err != nil {
		return "unknown"
	}
	return resp.Name
}

// Name implements plugins.Plugin interface for compatibility if needed.
func (c *GRPCClient) Name() string {
	return c.GetName()
}

// Version implements plugins.Plugin interface.
func (c *GRPCClient) Version() string {
	// Protocol version or plugin version?
	// Currently not exposed via RPC, hardcoded or needs RPC update.
	// For now return hardcoded as in previous implementation/mock.
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

func (c *GRPCClient) ConvertSchema(schema *models.Schema) (*models.Schema, error) {
	// Serialize domain object to proto string (JSON)
	schemaStr, err := ToProtoSchema(schema)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.ConvertSchema(context.Background(), &proto.ConvertSchemaRequest{Schema: schemaStr})
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	// Deserialize proto string (JSON) back to domain object
	return FromProtoSchema(resp.ConvertedSchema)
}
