package conversion

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

// Service is the interface for the conversion service.
type Service interface {
	ConvertFromFile(ctx context.Context, tracePath, outYaml string) error
}

// DefaultService is the default implementation of the conversion service.
type DefaultService struct {
	templateSvc *services.TemplateService
}

// NewService creates a new DefaultService.
func NewService() Service {
	return &DefaultService{templateSvc: services.NewTemplateService()}
}

// ConvertFromFile reads SQL traces from a file, converts them to templates, and saves them to a YAML file.
func (s *DefaultService) ConvertFromFile(ctx context.Context, tracePath, outYaml string) error {
	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var traces []models.SQLTrace
	dec := json.NewDecoder(file)
	for dec.More() {
		var t models.SQLTrace
		if err := dec.Decode(&t); err != nil {
			return err
		}
		traces = append(traces, t)
	}

	tc := models.TraceCollection{Traces: traces}
	tpls := s.templateSvc.ExtractTemplates(tc)

	f, err := os.Create(outYaml)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(map[string]interface{}{"templates": tpls})
}