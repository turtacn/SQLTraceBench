package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a benchmark from a workload file",
		RunE:  runRun,
	}
	runWorkloadPath string
	metricsPath     string
	runDB           string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runWorkloadPath, "workload", "w", "workload.json", "Path to the workload file")
	runCmd.Flags().StringVarP(&metricsPath, "out", "o", "metrics.json", "Path to the output metrics file")
	runCmd.Flags().StringVar(&runDB, "db", "", "Target database plugin to use (overrides config)")
}

func runRun(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()

	// Read the workload file.
	file, err := os.Open(runWorkloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var wl models.BenchmarkWorkload
	if err := json.NewDecoder(file).Decode(&wl); err != nil {
		return err
	}

	targetDB := cfg.Database.Driver
	if runDB != "" {
		targetDB = runDB
	}

	config := execution.ExecutionConfig{
		TargetDB:    targetDB,
		TargetQPS:   cfg.Benchmark.QPS,
		Concurrency: cfg.Benchmark.Concurrency,
	}

	metrics, err := root.Execution.RunBenchmark(context.Background(), &wl, config)
	if err != nil {
		return err
	}

	file, err = os.Create(metricsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metrics)
}