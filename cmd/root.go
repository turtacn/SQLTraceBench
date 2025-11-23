package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/pkg/config"
	"github.com/turtacn/SQLTraceBench/pkg/types"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

var (
	Version   = types.Version
	cfgFile   string
	pluginDir string
	cfg       *types.Config
	rootCmd   = &cobra.Command{
		Use:     "sqltracebench",
		Short:   "SQL trace-based workload benchmark CLI",
		Version: Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load the configuration.
			var err error
			cfg, err = config.Load(cfgFile)
			if err != nil {
				return err
			}

			// Initialize the logger.
			logger := utils.NewLogger(cfg.Log.Level, cfg.Log.Format, nil)
			utils.SetGlobalLogger(logger)

			// Load plugins
			if err := loadPlugins(); err != nil {
				return err
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cleanupPlugins()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", types.DefaultConfigPath, "config file")
	rootCmd.PersistentFlags().StringVar(&pluginDir, "plugin-dir", "./bin", "Directory where plugins are located")
}

func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd returns the root command for testing purposes.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// loadPlugins loads plugins from the configured directory.
func loadPlugins() error {
	return plugin_registry.LoadPlugins(pluginDir)
}

func cleanupPlugins() {
	plugin_registry.ClosePlugins()
}
