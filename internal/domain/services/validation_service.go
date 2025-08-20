package services

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type ValidationService interface {
	Validate(ctx context.Context, original, synthetic *models.PerformanceMetrics) (*models.ValidationReport, error)
}

type DefaultValidationService struct{}

func NewValidationService() ValidationService { return &DefaultValidationService{} }

func (s *DefaultValidationService) Validate(_ context.Context, o, n *models.PerformanceMetrics) (*models.ValidationReport, error) {
	r := &models.ValidationReport{OriginalStats: *o, SyntheticStats: *n}
	r.Compare(*o, *n)
	return r, nil
}
