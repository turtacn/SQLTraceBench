package validation

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

// Service is the interface for the validation service.
type Service interface {
	Validate(ctx context.Context, base, cand *models.PerformanceMetrics, threshold float64) (*models.ValidationReport, error)
}

// DefaultService is the default implementation of the validation service.
type DefaultService struct {
	validationSvc *services.ValidationService
	log           *utils.Logger
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{
		validationSvc: services.NewValidationService(),
		log:           utils.GetGlobalLogger(),
	}
}

// Validate compares two performance metrics and saves the report to a file.
func (s *DefaultService) Validate(ctx context.Context, base, cand *models.PerformanceMetrics, threshold float64) (*models.ValidationReport, error) {
	report := s.validationSvc.Validate(base, cand, threshold)
	s.log.Info("validation done", utils.Field{Key: "passed", Value: report.Pass})

	out, err := os.Create("validation.json")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if err := json.NewEncoder(out).Encode(report); err != nil {
		return nil, err
	}

	return report, nil
}