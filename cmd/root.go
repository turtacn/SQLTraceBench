package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

var (
	Version = types.Version
	cfgFile = types.DefaultConfigPath
	verbose bool
	rootCmd = &cobra.Command{
		Use:     "sqltracebench",
		Short:   "SQL trace-based workload benchmark CLI",
		Version: Version,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file (default is configs/default.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
}

func Execute() error {
	return rootCmd.Execute()
}
