package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert trace file to YAML templates",
	}
	convertTraceFile string
	convertOut       string
)

func init() {
	convertCmd.Flags().StringVarP(&convertTraceFile, "trace", "t", "trace.json", "trace JSON input")
	convertCmd.Flags().StringVarP(&convertOut, "out", "o", "templates.yaml", "YAML output file")
	rootCmd.AddCommand(convertCmd)
}
