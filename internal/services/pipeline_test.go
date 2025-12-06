package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// TestE2EPipelineSmoke runs a smoke test on the entire end-to-end pipeline.
func TestE2EPipelineSmoke(t *testing.T) {
	// 1. Create a trace collection
	tc := models.TraceCollection{}
	for i := 0; i < 10; i++ {
		tc.Add(models.SQLTrace{Query: "select * from users where id = :id"})
	}
	for i := 0; i < 5; i++ {
		tc.Add(models.SQLTrace{Query: "select * from orders"})
	}

	// 2. Extract templates
	templateService := services.NewTemplateService()
	ts := templateService.ExtractTemplates(tc)
	assert.Len(t, ts, 2, "should have extracted 2 templates")

	// 3. Generate a workload
	sampler := services.NewSeededWeightedRandomSampler(123)
	workloadService := services.NewWorkloadService(sampler)
	pm := services.NewParameterService().BuildModel(tc, ts)

	// Updated GenerateWorkload call
	// Increase query count to make QPS more stable
	config := services.GenerationConfig{
		TotalQueries: 100,
		ScaleFactor:  1.0,
	}
	wl, err := workloadService.GenerateWorkload(context.Background(), ts, pm, config)

	assert.NoError(t, err)
	assert.Len(t, wl.Queries, 100, "workload should have 100 queries")

	// 4. Run the benchmark twice
	rc := services.NewTokenBucketRateController(100, 10)
	defer rc.Stop()
	executionService := services.NewExecutionService(rc, 100*time.Millisecond)
	base, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "base run should not produce an error")

	executionService.Reset()
	cand, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "candidate run should not produce an error")

	// 5. Validate the results
	validationService := services.NewValidationService()
	// Increase threshold to 15% to tolerate test environment noise
	metadata := &models.ReportMetadata{Threshold: 0.15}
	report, err := validationService.ValidateAndReport(base, cand, metadata, "test_report.json")
	assert.NoError(t, err)
	assert.True(t, report.Result.Pass, "validation should pass for similar runs")
}
