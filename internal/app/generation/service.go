package generation

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// Service is the interface for the workload generation service.
type Service interface {
	GenerateWorkload(ctx context.Context, templates []models.SQLTemplate, n int) (*models.BenchmarkWorkload, error)
}

// DefaultService is the default implementation of the workload generation service.
type DefaultService struct {
	workloadSvc *services.WorkloadService
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{workloadSvc: services.NewWorkloadService()}
}

// GenerateWorkload creates a benchmark workload from a set of SQL templates.
func (s *DefaultService) GenerateWorkload(ctx context.Context, templates []models.SQLTemplate, n int) (*models.BenchmarkWorkload, error) {
	wl := s.workloadSvc.GenerateWorkload(templates, n)
	return &wl, nil
}