package services

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type WorkloadService interface {
	GenerateWorkload(ctx context.Context, tpls []*models.SQLTemplate, pm *models.ParameterModel) *models.BenchmarkWorkload
}

type DefaultWorkloadService struct {
	paramSvc ParameterService
}

func NewWorkloadService(ps ParameterService) WorkloadService {
	return &DefaultWorkloadService{paramSvc: ps}
}

func (s *DefaultWorkloadService) GenerateWorkload(
	_ context.Context,
	tpls []*models.SQLTemplate,
	pm *models.ParameterModel,
) *models.BenchmarkWorkload {
	w := &models.BenchmarkWorkload{
		Queries: make([]models.WorkloadQuery, 0, len(tpls)),
		Config:  models.ExecutionConfig{TargetQPS: 100},
	}
	for _, tpl := range tpls {
		// one concrete query per template
		params := s.paramSvc.GenerateParams(pm, tpl)
		sql, _ := tpl.GenerateQuery(params)
		w.Queries = append(w.Queries, models.WorkloadQuery{
			ID:     tpl.TemplateID,
			SQL:    sql,
			Weight: float64(tpl.Frequency),
		})
	}
	return w
}
