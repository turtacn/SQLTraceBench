// Package services contains the interfaces for the application's core services.
package services

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ValidationService is responsible for comparing the performance metrics of two benchmark runs.
type ValidationService struct{}

// NewValidationService creates a new ValidationService.
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// Validate compares the performance metrics of a base and candidate benchmark run.
// It returns a ValidationReport indicating whether the candidate's performance is within the acceptable threshold.
func (s *ValidationService) Validate(base, cand *models.PerformanceMetrics, th float64) *models.ValidationReport {
	vr := &models.ValidationReport{
		BaseQPS:      base.QPS(),
		CandidateQPS: cand.QPS(),
		DiffQPS:      cand.QPS() - base.QPS(),
		Threshold:    th,
	}

	// The candidate passes if its QPS is greater than or equal to the base QPS minus the threshold.
	vr.Pass = vr.CandidateQPS >= vr.BaseQPS*(1-th)

	return vr
}