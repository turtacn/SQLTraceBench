package plugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/turtacn/SQLTraceBench/pkg/plugin/grpc_impl"
)

// HandshakeConfig is the handshake config used by the host and the plugin.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SQLTRACEBENCH_PLUGIN",
	MagicCookieValue: "ON_FIRE",
}

// PluginMap is the map of plugins that we can serve.
var PluginMap = map[string]plugin.Plugin{
	"database_grpc": &grpc_impl.GRPCPluginImpl{},
}
