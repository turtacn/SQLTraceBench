// Package services contains the interfaces for the application's core services.
package services

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// WorkloadService is responsible for generating a benchmark workload from a set of SQL templates.
type WorkloadService struct{}

// NewWorkloadService creates a new WorkloadService.
func NewWorkloadService() *WorkloadService {
	return &WorkloadService{}
}

// GenerateWorkload creates a BenchmarkWorkload from a slice of SQL templates.
// For each template, it generates `n` instances of the query and adds them to the workload.
// In this version, parameter values are not yet randomized.
func (s *WorkloadService) GenerateWorkload(ts []models.SQLTemplate, n int) models.BenchmarkWorkload {
	var wl models.BenchmarkWorkload
	for _, t := range ts {
		// Create a dummy parameter map for query generation.
		// In a future implementation, this would involve parameter value generation.
		dummyParams := make(map[string]interface{})
		for _, p := range t.Parameters {
			// Use default values for now.
			dummyParams[p] = "default"
		}

		for i := 0; i < n; i++ {
			// In a more advanced implementation, we would generate different parameter values here.
			// For now, we just generate the same query n times.
			query, err := t.GenerateQuery(dummyParams)
			if err != nil {
				// For this simplified version, we'll ignore errors.
				// In a real application, this should be handled properly.
				continue
			}
			wl.Queries = append(wl.Queries, query)
		}
	}
	return wl
}