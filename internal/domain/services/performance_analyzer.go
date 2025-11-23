package services

import (
    "runtime"
    "sort"
    "github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type PerformanceAnalyzer struct{}

func (a *PerformanceAnalyzer) CalculateAverage(values []float64) float64 {
    if len(values) == 0 {
        return 0
    }
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

func (a *PerformanceAnalyzer) CalculatePercentile(values []float64, percentile float64) float64 {
    if len(values) == 0 {
        return 0
    }

    sorted := make([]float64, len(values))
    copy(sorted, values)
    sort.Float64s(sorted)

    index := int(float64(len(sorted)) * percentile)
    if index >= len(sorted) {
        index = len(sorted) - 1
    }

    return sorted[index]
}

func (a *PerformanceAnalyzer) GetResourceUsage() (memoryMB float64, cpuPercent float64) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    memoryMB = float64(m.Alloc) / 1024 / 1024

    // CPU usage is set to 0 as placeholder.
    cpuPercent = 0.0
    return
}

func (a *PerformanceAnalyzer) CalculateValidationScore(
    original, generated []models.SQLTrace,
    validator *StatisticalValidator,
) float64 {
    // Extract features for comparison
    origDuration := extractFeature(original, "duration")
    genDuration := extractFeature(generated, "duration")

    origTimestamp := extractFeature(original, "timestamp")
    genTimestamp := extractFeature(generated, "timestamp")

    // Run KS tests
    distResult := validator.KolmogorovSmirnovTest(origDuration, genDuration)
    temporalResult := validator.KolmogorovSmirnovTest(origTimestamp, genTimestamp)

    passCount := 0
    if distResult.Passed {
        passCount++
    }
    if temporalResult.Passed {
        passCount++
    }

    // Simple score calculation based on pass rate of these two tests
    return float64(passCount) / 2.0
}

func extractFeature(traces []models.SQLTrace, feature string) []float64 {
    result := make([]float64, len(traces))
    for i, trace := range traces {
        switch feature {
        case "duration":
            result[i] = float64(trace.Latency)
        case "timestamp":
            result[i] = float64(trace.Timestamp.Unix())
        }
    }
    return result
}
