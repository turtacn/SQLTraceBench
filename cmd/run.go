package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a benchmark from a workload file",
		RunE:  runRun,
	}
	runWorkloadPath string
	metricsPath     string
	qps             int
	concurrency     int
	slowThreshold   time.Duration
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runWorkloadPath, "workload", "w", "workload.json", "Path to the workload file")
	runCmd.Flags().StringVarP(&metricsPath, "out", "o", "metrics.json", "Path to the output metrics file")
	runCmd.Flags().IntVar(&qps, "qps", 100, "Target queries per second")
	runCmd.Flags().IntVar(&concurrency, "concurrency", 10, "Maximum concurrent queries")
	runCmd.Flags().DurationVar(&slowThreshold, "slow-threshold", 100*time.Millisecond, "Slow query threshold")
}

func runRun(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()
	metrics, err := root.Execution.RunBench(context.Background(), runWorkloadPath, qps, concurrency, slowThreshold)
	if err != nil {
		return err
	}

	file, err := os.Create(metricsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metrics)
}