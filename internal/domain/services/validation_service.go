package services

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/reports"
)

// ValidationService is responsible for comparing benchmark runs and generating a structured report.
type ValidationService struct {
	generator *reports.Generator
}

// NewValidationService creates a new ValidationService.
func NewValidationService() *ValidationService {
	return &ValidationService{
		generator: reports.NewGenerator(),
	}
}

// ValidateAndReport compares the performance of a base and candidate run, then generates and saves JSON and HTML reports.
func (s *ValidationService) ValidateAndReport(
	baseMetrics, candMetrics *models.PerformanceMetrics,
	metadata *models.ReportMetadata,
	outputPath string, // This serves as the prefix for the output files
) (*models.Report, error) {
	return s.generator.CompareAndGenerateReports(baseMetrics, candMetrics, metadata, outputPath)
}
