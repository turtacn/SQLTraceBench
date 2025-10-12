package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate the performance of a candidate benchmark against a base benchmark",
		RunE:  runValidate,
	}
	baseMetricsPath string
	candMetricsPath string
	reportPath      string
	threshold       float64
)

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVar(&baseMetricsPath, "base", "base_metrics.json", "Path to the base metrics file")
	validateCmd.Flags().StringVar(&candMetricsPath, "candidate", "candidate_metrics.json", "Path to the candidate metrics file")
	validateCmd.Flags().StringVarP(&reportPath, "out", "o", "report.json", "Path to the output report file")
	validateCmd.Flags().Float64Var(&threshold, "threshold", 0.05, "Performance degradation threshold")
}

func runValidate(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()
	_, err := root.Validation.Validate(context.Background(), baseMetricsPath, candMetricsPath, reportPath, threshold)
	return err
}