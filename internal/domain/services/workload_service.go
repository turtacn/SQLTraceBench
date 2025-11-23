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
	pm *models.WorkloadParameterModel,
	n int,
) (*models.BenchmarkWorkload, error) {
	var wl models.BenchmarkWorkload

	for _, t := range templates {
		if len(t.Parameters) == 0 {
			for i := 0; i < n; i++ {
				wl.Queries = append(wl.Queries, models.QueryWithArgs{
					Query: t.RawSQL,
					Args:  []interface{}{},
				})
			}
			continue
		}

		templateParams, ok := pm.TemplateParameters[t.GroupKey]
		if !ok {
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
					if len(dist.TopValues) > 0 {
						params[paramName] = dist.TopValues[0]
					} else {
						params[paramName] = "default"
					}
					continue
				}

				params[paramName] = sampledValue
			}

			queryWithArgs, err := t.GenerateQuery(params)
			if err != nil {
				continue
			}
			wl.Queries = append(wl.Queries, queryWithArgs)
		}
	}

	return &wl, nil
}
