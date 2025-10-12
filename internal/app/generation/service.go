package generation

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// Service is the interface for the workload generation service.
type Service interface {
	GenerateWorkload(ctx context.Context, tplPath string, n int) (*models.BenchmarkWorkload, error)
}

// DefaultService is the default implementation of the workload generation service.
type DefaultService struct {
	workloadSvc *services.WorkloadService
	paramSvc    *services.ParameterService
}

// NewService creates a new DefaultService.
func NewService() Service {
	sampler := services.NewWeightedRandomSampler()
	return &DefaultService{
		workloadSvc: services.NewWorkloadService(sampler),
		paramSvc:    services.NewParameterService(),
	}
}

// GenerateWorkload creates a benchmark workload from a set of SQL templates.
func (s *DefaultService) GenerateWorkload(ctx context.Context, tplPath string, n int) (*models.BenchmarkWorkload, error) {
	file, err := os.Open(tplPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var templates []models.SQLTemplate
	if err := json.NewDecoder(file).Decode(&templates); err != nil {
		return nil, err
	}

	// For now, we'll build the parameter model from the templates themselves, not a trace file.
	// This is a simplification for Phase 1.
	pm := s.paramSvc.BuildModel(models.TraceCollection{}, templates)

	wl, err := s.workloadSvc.GenerateWorkload(templates, pm, n)
	if err != nil {
		return nil, err
	}

	return wl, nil
}