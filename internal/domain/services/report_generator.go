package services

import (
	"fmt"
	"time"
)

// ValidationReport contains the aggregated results of all validation tests.
type ValidationReport struct {
	GeneratedAt       time.Time
	OverallScore      float64 // 0-100
	DistributionTests []ValidationResult
	TemporalTests     []ValidationResult
	QueryTypeTests    []ValidationResult
	Summary           string
}

// ReportGenerator generates validation reports from test results.
type ReportGenerator struct{}

// NewReportGenerator creates a new ReportGenerator.
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

// Generate aggregates validation results into a structured report.
func (g *ReportGenerator) Generate(
	distResults, temporalResults, queryResults []ValidationResult,
) *ValidationReport {
	totalTests := len(distResults) + len(temporalResults) + len(queryResults)
	passedTests := countPassed(distResults) + countPassed(temporalResults) + countPassed(queryResults)

	var overallScore float64
	if totalTests > 0 {
		overallScore = float64(passedTests) / float64(totalTests) * 100
	}

	summary := fmt.Sprintf(
		"Validation completed: %d/%d tests passed (%.1f%%)",
		passedTests, totalTests, overallScore,
	)

	return &ValidationReport{
		GeneratedAt:       time.Now(),
		OverallScore:      overallScore,
		DistributionTests: distResults,
		TemporalTests:     temporalResults,
		QueryTypeTests:    queryResults,
		Summary:           summary,
	}
}

func countPassed(results []ValidationResult) int {
	count := 0
	for _, r := range results {
		if r.Passed {
			count++
		}
	}
	return count
}
