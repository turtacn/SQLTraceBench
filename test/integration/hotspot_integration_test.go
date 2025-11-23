package integration

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/turtacn/SQLTraceBench/internal/domain/services"
    "github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestEndToEndAnalysis(t *testing.T) {
    // 1. Create a synthetic workload with Zipfian distribution
    traces := make([]models.SQLTrace, 0)

    // Simulate: param "id" follows Zipf (1 is very frequent)
    for i := 0; i < 1000; i++ {
        val := "1"
        if i % 10 == 0 { val = "2" }
        if i % 100 == 0 { val = "3" }

        traces = append(traces, models.SQLTrace{
            Query: "SELECT * FROM t WHERE id = :id",
            Parameters: map[string]interface{}{
                "id": val,
            },
        })
    }

    // 2. Run Analyzer
    analyzer := services.NewParameterAnalyzer()
    stats := analyzer.Analyze(traces)

    // 3. Verify Model
    model, ok := stats["id"]
    assert.True(t, ok)
    assert.Equal(t, models.DistZipfian, model.DistType)
    assert.Greater(t, model.HotspotRatio, 0.8) // "1" appears ~900 times out of 1000

    // Check TopValues
    assert.Equal(t, "1", model.TopValues[0])
}
