package services

import (
	"encoding/json"
	"fmt"
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
	pass := candMetrics.QPS() >= baseMetrics.QPS()*(1-metadata.Threshold)

	var reason string
	qpsDiff := candMetrics.QPS() - baseMetrics.QPS()
	qpsDiffPercent := (qpsDiff / baseMetrics.QPS()) * 100

	if pass {
		reason = fmt.Sprintf(
			"Validation passed. Candidate QPS of %.2f is within the %.2f%% threshold of the base QPS of %.2f (difference of %.2f, %.2f%%).",
			candMetrics.QPS(),
			metadata.Threshold*100,
			baseMetrics.QPS(),
			qpsDiff,
			qpsDiffPercent,
		)
	} else {
		reason = fmt.Sprintf(
			"Validation failed. Candidate QPS of %.2f is below the %.2f%% threshold of the base QPS of %.2f (difference of %.2f, %.2f%%).",
			candMetrics.QPS(),
			metadata.Threshold*100,
			baseMetrics.QPS(),
			qpsDiff,
			qpsDiffPercent,
		)
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
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return nil, err
	}

	return report, nil
}