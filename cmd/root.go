package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/pkg/config"
	"github.com/turtacn/SQLTraceBench/pkg/types"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

var (
	Version = types.Version
	cfgFile string
	cfg     *types.Config
	rootCmd = &cobra.Command{
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
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", types.DefaultConfigPath, "config file")
}

func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd returns the root command for testing purposes.
func GetRootCmd() *cobra.Command {
	return rootCmd
}