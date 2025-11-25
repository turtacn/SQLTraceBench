package generation

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// GenerateRequest encapsulates parameters for workload generation.
type GenerateRequest struct {
	SourceTraces []models.SQLTrace
	Count        int
}

// Service is the interface for the workload generation service.
type Service interface {
	GenerateWorkload(ctx context.Context, req GenerateRequest) (*models.BenchmarkWorkload, error)
}

// DefaultService is the default implementation of the workload generation service.
type DefaultService struct {
	templateSvc *services.TemplateService
	analyzer    *services.ParameterAnalyzer
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{
		templateSvc: services.NewTemplateService(),
		analyzer:    services.NewParameterAnalyzer(),
	}
}

func (s *DefaultService) GenerateWorkload(ctx context.Context, req GenerateRequest) (*models.BenchmarkWorkload, error) {
	if len(req.SourceTraces) == 0 {
		return nil, fmt.Errorf("generation requires source traces")
	}

	// 1. Parse traces to extract SQL templates and raw parameter values.
	templates := s.templateSvc.ExtractTemplates(models.TraceCollection{Traces: req.SourceTraces})
	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates could be extracted from the traces")
	}

	// 2. Build the statistical model for parameters.
	paramModels := s.analyzer.Analyze(req.SourceTraces)
	workloadModel := &models.WorkloadParameterModel{
		TemplateParameters: make(map[string]map[string]*models.ParameterModel),
	}

	// This is a simplified mapping. A more complex logic might be needed
	// if parameters are shared across different template groups.
	for _, tmpl := range templates {
		if _, ok := workloadModel.TemplateParameters[tmpl.GroupKey]; !ok {
			workloadModel.TemplateParameters[tmpl.GroupKey] = make(map[string]*models.ParameterModel)
		}
		for _, paramName := range tmpl.Parameters {
			if model, exists := paramModels[paramName]; exists {
				workloadModel.TemplateParameters[tmpl.GroupKey][paramName] = model
			}
		}
	}

	// 3. Synthesize the new workload.
	synth := services.NewSynthesizer(workloadModel)
	workload := &models.BenchmarkWorkload{
		Queries: make([]models.QueryWithArgs, 0, req.Count),
	}

	// 4. Generate the requested number of queries.
	totalWeight := 0
	for _, t := range templates {
		totalWeight += t.Weight
	}

	for i := 0; i < req.Count; i++ {
		// Select a template based on its observed frequency in the source traces.
		var template *models.SQLTemplate
		r := rand.Intn(totalWeight)
		currentW := 0
		for idx := range templates {
			currentW += templates[idx].Weight
			if r < currentW {
				template = &templates[idx]
				break
			}
		}
		if template == nil && len(templates) > 0 {
			template = &templates[len(templates)-1]
		}

		if template == nil {
			continue
		}

		// Fill the template's parameters with values sampled from the statistical model.
		args, err := synth.FillParameters(template)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not fill parameters for template %s: %v\n", template.GroupKey, err)
			continue
		}

		// Append the synthesized query to the workload.
		workload.Queries = append(workload.Queries, models.QueryWithArgs{
			Query: template.RawSQL,
			Args:  args,
		})
	}

	return workload, nil
}
