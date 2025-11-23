package services

import (
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// WorkloadService is responsible for generating a benchmark workload.
type WorkloadService struct {
	sampler Sampler
}

// NewWorkloadService creates a new WorkloadService.
func NewWorkloadService(sampler Sampler) *WorkloadService {
	return &WorkloadService{sampler: sampler}
}

// GenerateWorkload creates a BenchmarkWorkload from SQL templates and a parameter model.
func (s *WorkloadService) GenerateWorkload(
	templates []models.SQLTemplate,
	pm *models.ParameterModel,
	n int,
) (*models.BenchmarkWorkload, error) {
	var wl models.BenchmarkWorkload

	// We might want to distribute 'n' queries across templates based on their weights?
	// Currently it generates 'n' queries *per template*.
	// The CLI help says "Number of queries to generate per template".
	// So we keep this behavior.

	for _, t := range templates {
		// If the template has no parameters, we just generate the query as is.
		if len(t.Parameters) == 0 {
			for i := 0; i < n; i++ {
				// No params to replace
				wl.Queries = append(wl.Queries, models.QueryWithArgs{
					Query: t.RawSQL,
					Args:  []interface{}{},
				})
			}
			continue
		}

		templateParams, ok := pm.TemplateParameters[t.GroupKey]
		// If we don't have a model for this template, we can't generate parameterized queries reliably
		// unless we use defaults/dummies.
		if !ok {
			// Skip or warn? For now continue.
			continue
		}

		for i := 0; i < n; i++ {
			params := make(map[string]interface{})
			for _, paramName := range t.Parameters {
				dist, ok := templateParams[paramName]
				if !ok {
					params[paramName] = "default"
					continue
				}

				sampledValue, err := s.sampler.Sample(dist)
				if err != nil {
					// Fallback if sampling fails (e.g. empty distribution)
					// If distribution is empty, we should probably have a default value.
					params[paramName] = "default"
					continue
				}
				params[paramName] = sampledValue
			}

			// Generate the QueryWithArgs struct.
			queryWithArgs, err := t.GenerateQuery(params)
			if err != nil {
				continue
			}
			wl.Queries = append(wl.Queries, queryWithArgs)
		}
	}

	return &wl, nil
}
