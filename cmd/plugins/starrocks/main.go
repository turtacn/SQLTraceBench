package main

import (
	"github.com/hashicorp/go-plugin"
	pkgplugin "github.com/turtacn/SQLTraceBench/pkg/plugin"
	"github.com/turtacn/SQLTraceBench/pkg/plugin/grpc_impl"
	"github.com/turtacn/SQLTraceBench/plugins/starrocks"
)

func main() {
	// Initialize the real implementation
	impl := starrocks.New()

	// Serve using go-plugin's standardized mechanism
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pkgplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"database_plugin": &grpc_impl.GRPCPluginImpl{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
