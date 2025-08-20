package cmd

import (
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate synthetic workload against original trace",
	}
	validateFile string
)

func init() {
	validateCmd.Flags().StringVarP(&validateFile, "result", "r", "bench_out.json", "benchmark result file")
	rootCmd.AddCommand(validateCmd)
}
