// Package services contains the interfaces for the application's core services.
package services

import (
	"sort"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// TemplateService is responsible for processing raw SQL traces and extracting SQL templates.
// It normalizes SQL queries to group them into templates and calculates the frequency of each template.
type TemplateService struct{}

// NewTemplateService creates a new TemplateService.
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// ExtractTemplates takes a collection of SQL traces and returns a slice of SQL templates.
// It groups traces by a normalized version of their SQL query, calculates the weight (frequency) of each template,
// and sorts the templates by weight in descending order.
func (s *TemplateService) ExtractTemplates(tc models.TraceCollection) []models.SQLTemplate {
	agg := make(map[string]*models.SQLTemplate)

	for _, tr := range tc.Traces {
		// Normalize the query to use as a grouping key.
		// This is a simple normalization; more sophisticated techniques could be used here.
		key := strings.ToLower(strings.TrimSpace(tr.Query))

		if _, ok := agg[key]; !ok {
			template := &models.SQLTemplate{
				RawSQL:   tr.Query,
				GroupKey: key,
				Weight:   0,
			}
			template.ExtractParameters()
			agg[key] = template
		}
		agg[key].Weight++
	}

	// Convert the map to a slice and sort it by weight.
	out := make([]models.SQLTemplate, 0, len(agg))
	for _, t := range agg {
		out = append(out, *t)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Weight > out[j].Weight
	})

	return out
}