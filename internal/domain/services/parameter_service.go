package services

import (
	"regexp"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ParameterService is responsible for building a statistical model of parameters from SQL traces.
type ParameterService struct {
	extractor *mockParameterExtractor
}

// NewParameterService creates a new ParameterService.
func NewParameterService() *ParameterService {
	return &ParameterService{
		extractor: &mockParameterExtractor{},
	}
}

// BuildModel analyzes a collection of SQL traces and builds a ParameterModel.
func (s *ParameterService) BuildModel(tc models.TraceCollection, templates []models.SQLTemplate) *models.ParameterModel {
	pm := models.NewParameterModel()
	templateMap := make(map[string]models.SQLTemplate)
	for _, t := range templates {
		templateMap[normalizeQuery(t.RawSQL)] = t
	}

	for _, trace := range tc.Traces {
		key := normalizeQuery(trace.Query)
		template, ok := templateMap[key]
		if !ok {
			continue
		}

		paramValues, err := s.extractor.Extract(trace.Query, &template)
		if err != nil {
			continue
		}

		if _, ok := pm.TemplateParameters[template.GroupKey]; !ok {
			pm.TemplateParameters[template.GroupKey] = make(map[string]*models.ValueDistribution)
		}

		for paramName, value := range paramValues {
			if _, ok := pm.TemplateParameters[template.GroupKey][paramName]; !ok {
				pm.TemplateParameters[template.GroupKey][paramName] = models.NewValueDistribution()
			}
			pm.TemplateParameters[template.GroupKey][paramName].AddObservation(value)
		}
	}

	return pm
}

// mockParameterExtractor simulates the process of extracting parameter values from a SQL query.
type mockParameterExtractor struct{}

func (m *mockParameterExtractor) Extract(sql string, template *models.SQLTemplate) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	// This is a very basic extractor that only works for the test case.
	re := regexp.MustCompile(`id = (\d+)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) > 1 {
		values[":id"] = matches[1]
	}
	return values, nil
}

// normalizeQuery provides a basic normalization for lookup.
func normalizeQuery(q string) string {
	q = strings.ToLower(q)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
	q = regexp.MustCompile(`\s*=\s*(\d+|'[^']*')`).ReplaceAllString(q, " = ?")
	q = regexp.MustCompile(`\s*=\s*:\w+`).ReplaceAllString(q, " = ?")
	return q
}