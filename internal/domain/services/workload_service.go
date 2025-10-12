package services

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// WorkloadService is responsible for generating a benchmark workload.
// It uses a Sampler to generate realistic parameter values based on a ParameterModel.
type WorkloadService struct {
	sampler Sampler
}

// NewWorkloadService creates a new WorkloadService.
func NewWorkloadService(sampler Sampler) *WorkloadService {
	return &WorkloadService{sampler: sampler}
}

// GenerateWorkload creates a BenchmarkWorkload from SQL templates and a parameter model.
// It generates `n` queries for each template, using the sampler to select parameter values
// based on their observed distribution.
func (s *WorkloadService) GenerateWorkload(
	templates []models.SQLTemplate,
	pm *models.ParameterModel,
	n int,
) (*models.BenchmarkWorkload, error) {
	var wl models.BenchmarkWorkload

	for _, t := range templates {
		templateParams, ok := pm.TemplateParameters[t.GroupKey]
		if !ok {
			// If there's no model for this template, we can't generate queries.
			// In a real application, we might want to handle this differently (e.g., use default values).
			continue
		}

		for i := 0; i < n; i++ {
			params := make(map[string]interface{})
			for _, paramName := range t.Parameters {
				dist, ok := templateParams[paramName]
				if !ok {
					// No distribution found for this parameter.
					// We could use a default value or return an error.
					params[paramName] = "default" // Fallback
					continue
				}

				// Sample a value from the distribution.
				sampledValue, err := s.sampler.Sample(dist)
				if err != nil {
					// Handle the error, e.g., by logging it and using a default.
					params[paramName] = "default"
					continue
				}
				params[paramName] = sampledValue
			}

			// Generate the final query with the sampled parameters.
			query, err := t.GenerateQuery(params)
			if err != nil {
				// Log or handle the error.
				continue
			}
			wl.Queries = append(wl.Queries, query)
		}
	}

	return &wl, nil
}