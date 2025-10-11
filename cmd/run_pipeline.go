package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/turtacn/SQLTraceBench/internal/app"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
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

// RunPipeline executes a demonstration of the end-to-end benchmark pipeline.
// It uses in-memory data and simulated services to showcase the workflow.
func RunPipeline(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log := utils.GetGlobalLogger()
	root := app.NewRoot()

	// Step 1: Create a dummy trace collection
	log.Info("Step 1: Creating trace collection...")
	tc := models.TraceCollection{}
	for i := 0; i < 10; i++ {
		tc.Add(models.SQLTrace{Query: "select * from users where id = :id"})
	}
	for i := 0; i < 5; i++ {
		tc.Add(models.SQLTrace{Query: "select * from orders"})
	}
	// Note: In a real scenario, this would come from a file via the conversion service.
	// For this command, we'll just extract the templates directly.
	templateService := &services.TemplateService{}
	templates := templateService.ExtractTemplates(tc)
	log.Info("Extracted templates", utils.Field{Key: "count", Value: len(templates)})

	// Step 2: Generate workload
	log.Info("Step 2: Generating workload...")
	workload, err := root.Generation.GenerateWorkload(ctx, templates, 10)
	if err != nil {
		log.Error("generation step failed", utils.Field{Key: "error", Value: err})
		return err
	}
	log.Info("Generated workload", utils.Field{Key: "query_count", Value: len(workload.Queries)})

	// Step 3: Run benchmark twice (base and candidate)
	log.Info("Step 3: Running benchmark...")
	baseMetrics, err := root.Execution.RunBench(ctx, workload)
	if err != nil {
		log.Error("base run failed", utils.Field{Key: "error", Value: err})
		return err
	}
	log.Info("Base run complete", utils.Field{Key: "qps", Value: baseMetrics.QPS()})

	candidateMetrics, err := root.Execution.RunBench(ctx, workload)
	if err != nil {
		log.Error("candidate run failed", utils.Field{Key: "error", Value: err})
		return err
	}
	log.Info("Candidate run complete", utils.Field{Key: "qps", Value: candidateMetrics.QPS()})

	// Step 4: Validate the results
	log.Info("Step 4: Validating...")
	report, err := root.Validation.Validate(ctx, baseMetrics, candidateMetrics, 0.05)
	if err != nil {
		log.Error("validation step failed", utils.Field{Key: "error", Value: err})
		return err
	}

	log.Info("Pipeline completed.", utils.Field{Key: "passed", Value: report.Pass})
	return nil
}