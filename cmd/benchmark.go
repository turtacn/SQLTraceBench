package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app/benchmark"
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run performance benchmarks",
}

var benchmarkRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run benchmark tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		outputDir, _ := cmd.Flags().GetString("output")
		exportProm, _ := cmd.Flags().GetBool("prometheus")

		service := benchmark.NewDefaultBenchmarkService()

		if exportProm {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				fmt.Println("Starting Prometheus metrics server on :9091")
				if err := http.ListenAndServe(":9091", nil); err != nil {
					fmt.Printf("Error starting Prometheus server: %v\n", err)
				}
			}()
		}

		report, err := service.RunBenchmark(context.Background(), benchmark.BenchmarkRequest{
			ConfigPath:       configPath,
			OutputDir:        outputDir,
			ExportPrometheus: exportProm,
		})

		if err != nil {
			return fmt.Errorf("benchmark failed: %w", err)
		}

		fmt.Printf("Benchmark completed!\n")
		fmt.Printf("Summary: %s\n", report.Summary)
		fmt.Printf("Report: %s\n", report.ReportPath)

		if exportProm {
			fmt.Println("Prometheus metrics are exposed at :9091/metrics")
			fmt.Println("Press Ctrl+C to exit...")
			select {} // Block forever
		}

		return nil
	},
}

func init() {
	benchmarkCmd.AddCommand(benchmarkRunCmd)

	benchmarkRunCmd.Flags().String("config", "configs/benchmark.yaml", "Benchmark config file")
	benchmarkRunCmd.Flags().String("output", "./benchmark_reports", "Output directory")
	benchmarkRunCmd.Flags().Bool("prometheus", true, "Export Prometheus metrics")

	rootCmd.AddCommand(benchmarkCmd)
}
