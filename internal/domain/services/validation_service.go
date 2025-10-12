package services

import (
	"encoding/json"
	"os"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ValidationService is responsible for comparing benchmark runs and generating a structured report.
type ValidationService struct{}

// NewValidationService creates a new ValidationService.
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateAndReport compares the performance of a base and candidate run, then generates and saves a JSON report.
func (s *ValidationService) ValidateAndReport(
	baseMetrics, candMetrics *models.PerformanceMetrics,
	metadata *models.ReportMetadata,
	outputPath string,
) (*models.Report, error) {

	// Perform the validation logic.
	// For this phase, we'll keep it simple: pass if the candidate's QPS is within the threshold.
	pass := candMetrics.QPS() >= baseMetrics.QPS()*(1-metadata.Threshold)
	reason := "Validation passed: Candidate QPS is within the acceptable threshold."
	if !pass {
		reason = "Validation failed: Candidate QPS is below the acceptable threshold."
	}

	// Assemble the report.
	report := &models.Report{
		Version:   "report.v1",
		Timestamp: time.Now(),
		Metadata:  metadata,
		Result: &models.ValidationResult{
			BaseMetrics:      baseMetrics,
			CandidateMetrics: candMetrics,
			Pass:             pass,
			Reason:           reason,
		},
	}

	// Save the report to a file.
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	if err := encoder.Encode(report); err != nil {
		return nil, err
	}

	return report, nil
}