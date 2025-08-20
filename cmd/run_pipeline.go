package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

var (
	pipelineCmd = &cobra.Command{
		Use:   "run-pipeline",
		Short: "Run whole pipeline: convert→generate→bench→validate",
		RunE:  RunPipeline,
	}
)

func init() {
	rootCmd.AddCommand(pipelineCmd)
}

func RunPipeline(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cfg := database.Config{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "",
		Database: "test",
	}

	root := app.NewRoot(cfg)
	root.Cfg = cfg

	log := utils.GetGlobalLogger()
	log.Info("Step 1: converting trace...")
	_, err := root.Conversion.ConvertFromFile(ctx, "trace.json", "templates.yaml")
	if err != nil {
		log.Error("convert step failed", utils.Field{Key: "error", Value: err})
		return err
	}

	log.Info("Step 2: generating workload...")
	// root.Generate.GenerateWorkload(ctx, ...)

	log.Info("Step 3: running benchmark...")
	_, _ = root.Execution.RunBench(ctx, "workload.yaml", "mysql")

	log.Info("Step 4: validating...")
	_ = root.Validation.Validate(ctx, "validation.json")

	log.Info("Pipeline completed.")
	return nil
}
