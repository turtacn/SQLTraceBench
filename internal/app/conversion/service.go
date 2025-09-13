package conversion

import (
	"context"
	"encoding/json"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

type Service interface {
	ConvertFromFile(ctx context.Context, tracePath, outYaml string) error
}

type DefaultService struct {
	templateSvc services.TemplateService
	schemaSvc   services.SchemaService
}

func NewService(ts services.TemplateService, ss services.SchemaService) Service {
	return &DefaultService{templateSvc: ts, schemaSvc: ss}
}

func (s *DefaultService) ConvertFromFile(ctx context.Context, tracePath, outYaml string) error {
	var traces []models.SQLTrace
	file, _ := os.Open(tracePath)
	defer file.Close()
	dec := json.NewDecoder(file)
	for dec.More() {
		var t models.SQLTrace
		_ = dec.Decode(&t)
		traces = append(traces, t)
	}

	tpls, _ := s.templateSvc.ExtractTemplates(ctx, traces)
	f, _ := os.Create(outYaml)
	defer f.Close()
	return json.NewEncoder(f).Encode(map[string]interface{}{"templates": tpls})
}
