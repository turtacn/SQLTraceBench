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
	ConvertFromFile(ctx context.Context, tracePath string) ([]models.SQLTemplate, error)
}

// DefaultService is the default implementation of the conversion service.
type DefaultService struct {
	templateSvc *services.TemplateService
	parser      services.Parser
}

// NewService creates a new DefaultService.
func NewService(parser services.Parser) Service {
	return &DefaultService{
		templateSvc: services.NewTemplateService(),
		parser:      parser,
	}
}

// ConvertFromFile reads SQL traces from a file and converts them to templates.
func (s *DefaultService) ConvertFromFile(ctx context.Context, tracePath string) ([]models.SQLTemplate, error) {
	file, err := os.Open(tracePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var traces []models.SQLTrace
	dec := json.NewDecoder(file)
	for dec.More() {
		var t models.SQLTrace
		if err := dec.Decode(&t); err != nil {
			return nil, err
		}
		traces = append(traces, t)
	}

	tc := models.TraceCollection{Traces: traces}
	tpls := s.templateSvc.ExtractTemplates(tc)

	// Use the parser to extract table names for each template.
	for i := range tpls {
		tables, err := s.parser.ListTables(tpls[i].RawSQL)
		if err != nil {
			// In a real application, we might want to handle this error more gracefully.
			continue
		}
		_ = tables // TODO: store the tables in the template
	}

	return tpls, nil
}