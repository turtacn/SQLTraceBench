package types

import "time"

// Config holds the application's configuration.
type Config struct {
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Benchmark BenchmarkConfig `mapstructure:"benchmark"`
}

// LogConfig holds the logging configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// BenchmarkConfig holds the benchmark configuration.
type BenchmarkConfig struct {
	Executor      string        `mapstructure:"executor"`
	QPS           int           `mapstructure:"qps"`
	Concurrency   int           `mapstructure:"concurrency"`
	SlowThreshold time.Duration `mapstructure:"slow_threshold"`
}

var DefaultConfigPath = "configs/default.yaml"
var Version = "dev"