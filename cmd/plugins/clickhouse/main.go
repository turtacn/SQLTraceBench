package main

import (
	"github.com/hashicorp/go-plugin"
	pkg_plugin "github.com/turtacn/SQLTraceBench/pkg/plugin"
	"github.com/turtacn/SQLTraceBench/pkg/plugin/grpc_impl"
	"github.com/turtacn/SQLTraceBench/plugins/clickhouse"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	impl := clickhouse.New() // Instantiate business logic

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pkg_plugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"database_plugin": &grpc_impl.GRPCPluginImpl{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
