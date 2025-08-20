package types

const (
	Version           = "1.0.0"
	AppName           = "SQLTraceBench"
	DefaultConfigPath = "./configs/default.yaml"
	DefaultTempDir    = "/tmp/sqltracebench"
	DefaultPluginDir  = "./plugins"

	MaxSQLLength       = 1024 * 1024 // 1MB
	DefaultTimeout     = 30 * 1e9    // 30s
	DefaultBatchSize   = 1000
	DefaultConcurrency = 4

	DefaultConnectionTimeout = 10 * 1e9
	DefaultQueryTimeout      = 30 * 1e9
	MaxConnections           = 100
	MaxRetries               = 3

	DefaultQPS       = 100.0
	MaxConcurrency   = 64
	DefaultDuration  = 5 * 1e9 * 60 // 5min
	DefaultReportDir = "./reports"
)

//Personal.AI order the ending
