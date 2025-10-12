package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
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
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&genTplPath, "templates", "t", "templates.json", "Path to the SQL template file")
	generateCmd.Flags().StringVarP(&workloadPath, "out", "o", "workload.json", "Path to the output workload file")
	generateCmd.Flags().IntVarP(&genCount, "count", "c", 10, "Number of queries to generate per template")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	root := app.NewRoot()
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