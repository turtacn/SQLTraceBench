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
	wl, err := workloadService.GenerateWorkload(ts, pm, 1)
	assert.NoError(t, err)
	assert.Len(t, wl.Queries, 2, "workload should have 2 queries")

	// 4. Run the benchmark twice
	rc := services.NewTokenBucketRateController(100, 10)
	defer rc.Stop()
	executionService := services.NewExecutionService(rc, 100*time.Millisecond)
	base, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "base run should not produce an error")
	cand, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "candidate run should not produce an error")

	// 5. Validate the results
	validationService := services.NewValidationService()
	metadata := &models.ReportMetadata{Threshold: 0.05}
	report, err := validationService.ValidateAndReport(base, cand, metadata, "test_report.json")
	assert.NoError(t, err)
	assert.True(t, report.Result.Pass, "validation should pass for similar runs")
}