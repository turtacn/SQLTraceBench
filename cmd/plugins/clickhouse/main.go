package main

import (
	"github.com/hashicorp/go-plugin"
	pkg_plugin "github.com/turtacn/SQLTraceBench/pkg/plugin"
	"github.com/turtacn/SQLTraceBench/pkg/plugin/grpc_impl"
	"github.com/turtacn/SQLTraceBench/plugins/clickhouse"
)

func main() {
	// Initialize the actual business logic
	// clickhouse.New() returns plugins.Plugin (impl in plugins/clickhouse/plugin.go)
	realImpl := clickhouse.New()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pkg_plugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"database_grpc": &grpc_impl.GRPCPluginImpl{Impl: realImpl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
