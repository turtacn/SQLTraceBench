package cmd

import (
	"context"
	"encoding/json"
	"os"

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
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runWorkloadPath, "workload", "w", "workload.json", "Path to the workload file")
	runCmd.Flags().StringVarP(&metricsPath, "out", "o", "metrics.json", "Path to the output metrics file")
}

func runRun(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()
	metrics, err := root.Execution.RunBench(
		context.Background(),
		runWorkloadPath,
		cfg.Benchmark.Executor,
		cfg.Database.Driver,
		cfg.Database.DSN,
		cfg.Benchmark.QPS,
		cfg.Benchmark.Concurrency,
		cfg.Benchmark.SlowThreshold,
	)
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