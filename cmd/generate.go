package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a workload file from a SQL template file",
		RunE:  runGenerate,
	}
	genTplPath   string
	workloadPath string
	genCount     int
	hotspotRatio float64
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&genTplPath, "templates", "t", "templates.json", "Path to the SQL template file")
	generateCmd.Flags().StringVarP(&workloadPath, "out", "o", "workload.json", "Path to the output workload file")
	generateCmd.Flags().IntVarP(&genCount, "count", "c", 10, "Number of queries to generate per template")
	generateCmd.Flags().Float64Var(&hotspotRatio, "hotspot-ratio", 0.0, "Zipf skew parameter for hotspot injection (e.g., 1.1). 0 to disable.")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()

	// Configure sampler based on flags
	if hotspotRatio > 0 {
		sampler := services.NewZipfSampler(hotspotRatio)
		root.Generation.SetSampler(sampler)
	}

	workload, err := root.Generation.GenerateWorkload(context.Background(), genTplPath, genCount)
	if err != nil {
		return err
	}

	file, err := os.Create(workloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(workload)
}
