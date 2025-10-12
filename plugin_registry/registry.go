package plugin_registry

import (
	"github.com/turtacn/SQLTraceBench/plugins"
	"github.com/turtacn/SQLTraceBench/plugins/clickhouse"
	"github.com/turtacn/SQLTraceBench/plugins/starrocks"
)

func init() {
	// Automatically register all available plugins.
	plugins.GlobalRegistry.Register(clickhouse.New())
	plugins.GlobalRegistry.Register(starrocks.New())
}