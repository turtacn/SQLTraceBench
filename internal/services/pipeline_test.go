package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// TestE2EPipelineSmoke runs a smoke test on the entire end-to-end pipeline.
// It simulates the full workflow:
// 1. A collection of SQL traces is created.
// 2. Templates are extracted from the traces.
// 3. A workload is generated from the templates.
// 4. The workload is executed twice to get base and candidate performance metrics.
// 5. The metrics are validated to ensure the results are as expected.
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
	workloadService := services.NewWorkloadService()
	wl := workloadService.GenerateWorkload(ts, 2)
	// (10 users * 2) + (5 orders * 2) = 30, but we generate n per template, so 2*2 = 4 queries
	// The number of queries should be 2 * len(ts)
	assert.Len(t, wl.Queries, 4, "workload should have 4 queries")

	// 4. Run the benchmark twice
	executionService := services.NewExecutionService()
	base, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "base run should not produce an error")
	cand, err := executionService.RunBench(context.Background(), wl)
	assert.NoError(t, err, "candidate run should not produce an error")

	// 5. Validate the results
	validationService := services.NewValidationService()
	vr := validationService.Validate(base, cand, 0.05)
	assert.True(t, vr.Pass, "validation should pass for similar runs")
	assert.InDelta(t, base.QPS(), cand.QPS(), 10, "QPS should be similar for identical workloads")
}