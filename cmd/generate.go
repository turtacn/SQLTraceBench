package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a workload file from a source SQL trace file",
		RunE:  runGenerate,
	}
	sourceTracePath string
	workloadPath    string
	genCount        int
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&sourceTracePath, "source-traces", "s", "traces.json", "Path to the source SQL trace file")
	generateCmd.Flags().StringVarP(&workloadPath, "out", "o", "workload.json", "Path to the output workload file")
	generateCmd.Flags().IntVarP(&genCount, "count", "c", 1000, "Number of queries to generate in the workload")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()

	// 1. Load source traces from the input file.
	traceData, err := ioutil.ReadFile(sourceTracePath)
	if err != nil {
		return err
	}
	var sourceTraces []models.SQLTrace
	if err := json.Unmarshal(traceData, &sourceTraces); err != nil {
		return err
	}

	// 2. Create the generation request.
	req := generation.GenerateRequest{
		SourceTraces: sourceTraces,
		Count:        genCount,
	}

	// 3. Generate the workload.
	workload, err := root.Generation.GenerateWorkload(context.Background(), req)
	if err != nil {
		return err
	}

	// 4. Write the workload to the output file.
	file, err := os.Create(workloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(workload)
}
