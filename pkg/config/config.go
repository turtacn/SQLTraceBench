package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// Load loads the application configuration from a file and environment variables.
func Load(path string) (*types.Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Set default values.
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("benchmark.executor", "simulated")
	v.SetDefault("benchmark.qps", 100)
	v.SetDefault("benchmark.concurrency", 10)
	v.SetDefault("benchmark.slow_threshold", "100ms")

	// Allow environment variables to override config file settings.
	v.SetEnvPrefix("SQLTRACEBENCH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// Ignore if the config file doesn't exist; we'll use the defaults.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg types.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}