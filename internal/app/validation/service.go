package validation

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// Service is the interface for the validation service.
type Service interface {
	Validate(ctx context.Context, basePath, candPath, outputPath string, threshold float64) (*models.Report, error)
}

// DefaultService is the default implementation of the validation service.
type DefaultService struct {
	validationSvc *services.ValidationService
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{
		validationSvc: services.NewValidationService(),
	}
}

// Validate compares two performance metrics from files and saves the report.
func (s *DefaultService) Validate(ctx context.Context, basePath, candPath, outputPath string, threshold float64) (*models.Report, error) {
	baseMetrics, err := readMetricsFile(basePath)
	if err != nil {
		return nil, err
	}

	candMetrics, err := readMetricsFile(candPath)
	if err != nil {
		return nil, err
	}

	metadata := &models.ReportMetadata{
		BaseTarget:      basePath,
		CandidateTarget: candPath,
		Threshold:       threshold,
	}

	report, err := s.validationSvc.ValidateAndReport(baseMetrics, candMetrics, metadata, outputPath)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func readMetricsFile(path string) (*models.PerformanceMetrics, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metrics models.PerformanceMetrics
	if err := json.NewDecoder(file).Decode(&metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}