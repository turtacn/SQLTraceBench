package main

import (
	"github.com/hashicorp/go-plugin"
	pkg_plugin "github.com/turtacn/SQLTraceBench/pkg/plugin"
	"github.com/turtacn/SQLTraceBench/pkg/plugin/grpc_impl"
	"github.com/turtacn/SQLTraceBench/plugins/starrocks"
)

func main() {
	// Create the plugin implementation
	impl := starrocks.New()

	// Serve the plugin using the standard configuration
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pkg_plugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"database_plugin": &grpc_impl.GRPCPluginImpl{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
