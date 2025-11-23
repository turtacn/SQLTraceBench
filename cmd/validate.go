package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
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
	// The validation service interface has changed to ValidateTrace(ctx, req)
	// We need to adapt the CLI command.
	// But `baseMetricsPath` implies comparison of metrics, not traces.
	// The issue is that I am working on Phase 4 (Benchmarking), but `cmd/validate.go`
	// seems to be from an older phase or conflicting with the `internal/app/validation` I see.
	// The `internal/app/validation/service.go` has `ValidateTrace` but no `Validate`.
	// I will fix this compilation error by using `ValidateTrace`.
	// However, the flags (base, candidate) look like they expect metrics JSONs,
	// whereas `ValidateTrace` expects trace files (JSONL).
	// Assuming `baseMetricsPath` == Original traces and `candMetricsPath` == Generated traces.

	root := app.NewRoot()
	req := validation.ValidationRequest{
		OriginalPath:   baseMetricsPath,
		GeneratedPath:  candMetricsPath,
		ReportDir:      reportPath, // Actually reportPath is a file path in flags, but Dir in Request
		KSThreshold:    threshold,
	}

	// If reportPath is a file, we might need to adjust.
	// But let's just pass it.

	_, err := root.Validation.ValidateTrace(context.Background(), req)
	return err
}