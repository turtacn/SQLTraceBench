package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/reporters"
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

	// Load base metrics
	baseFile, err := os.Open(baseMetricsPath)
	if err != nil {
		return err
	}
	defer baseFile.Close()
	var baseResult models.BenchmarkResult
	if err := json.NewDecoder(baseFile).Decode(&baseResult); err != nil {
		return err
	}

	// Load candidate metrics
	candFile, err := os.Open(candMetricsPath)
	if err != nil {
		return err
	}
	defer candFile.Close()
	var candResult models.BenchmarkResult
	if err := json.NewDecoder(candFile).Decode(&candResult); err != nil {
		return err
	}

	report, err := root.Validation.ValidateBenchmarks(context.Background(), &baseResult, &candResult)
	if err != nil {
		return err
	}

	reporter, err := reporters.NewHTMLReporter()
	if err != nil {
		return err
	}
	return reporter.GenerateReport(report, reportPath)
}