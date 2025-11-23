package services

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ParameterExtractor defines the interface for extracting parameters from SQL queries.
type ParameterExtractor interface {
	Extract(sql string, template *models.SQLTemplate) (map[string]interface{}, error)
}

// ParameterService is responsible for building a statistical model of parameters from SQL traces.
type ParameterService struct {
	extractor ParameterExtractor
	analyzer  *ParameterAnalyzer
}

// NewParameterService creates a new ParameterService.
func NewParameterService() *ParameterService {
	return &ParameterService{
		extractor: &RegexParameterExtractor{
			regexCache: make(map[string]*regexp.Regexp),
		},
		analyzer: NewParameterAnalyzer(),
	}
}

// BuildModel analyzes a collection of SQL traces and builds a WorkloadParameterModel.
func (s *ParameterService) BuildModel(tc models.TraceCollection, templates []models.SQLTemplate) *models.WorkloadParameterModel {
	pm := models.NewWorkloadParameterModel()

	// 1. Group traces by template
	// We need a way to match trace to template.
	// Assuming traces have their query. We match by normalizing.

	templateMap := make(map[string]*models.SQLTemplate)
	tracesByTemplate := make(map[string][]models.SQLTrace)

	for i := range templates {
		t := &templates[i]
		// Ensure parameters are extracted
		if len(t.Parameters) == 0 {
			t.ExtractParameters()
		}

		key := normalizeQuery(t.RawSQL)
		templateMap[key] = t

		// Initialize the map entry for this template
		if _, ok := pm.TemplateParameters[t.GroupKey]; !ok {
			pm.TemplateParameters[t.GroupKey] = make(map[string]*models.ParameterModel)
		}
	}

	for _, trace := range tc.Traces {
		key := normalizeQuery(trace.Query)
		template, ok := templateMap[key]
		if !ok {
			// Try to find best match? Or skip.
			continue
		}

		// Extract parameters
		paramValues, err := s.extractor.Extract(trace.Query, template)
		if err != nil {
			continue
		}

		// Enrich trace with parameters
		// We make a copy of trace to avoid modifying original if needed, or just use a new struct
		t := trace
		t.Parameters = paramValues
		tracesByTemplate[template.GroupKey] = append(tracesByTemplate[template.GroupKey], t)
	}

	// 2. Analyze parameters for each template
	for groupKey, traces := range tracesByTemplate {
		// Use ParameterAnalyzer
		paramModels := s.analyzer.Analyze(traces)
		pm.TemplateParameters[groupKey] = paramModels
	}

	// 3. Handle templates with no traces (Default/Uniform)
	for _, t := range templates {
		if len(tracesByTemplate[t.GroupKey]) == 0 {
			// Create default models for parameters
			for _, pName := range t.Parameters {
				pm.TemplateParameters[t.GroupKey][pName] = &models.ParameterModel{
					ParamName:    pName,
					DataType:     "STRING", // default
					DistType:     models.DistUniform,
					Cardinality:  1,
					TopValues:    []interface{}{"1"},
					TopFrequencies: []int{1},
				}
			}
		}
	}

	return pm
}

// RegexParameterExtractor extracts parameter values using regex matching against the template.
type RegexParameterExtractor struct {
	regexCache map[string]*regexp.Regexp
	mu         sync.RWMutex
}

func (r *RegexParameterExtractor) Extract(sql string, template *models.SQLTemplate) (map[string]interface{}, error) {
	// We need to convert the template RawSQL into a regex that captures values.
	// Template: SELECT * FROM users WHERE id = :id AND name = :name
	// Regex:    SELECT * FROM users WHERE id = (.*?) AND name = (.*?)

	r.mu.RLock()
	re, ok := r.regexCache[template.RawSQL]
	r.mu.RUnlock()

	if !ok {
		// Re-using the regex logic
		paramRe := regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)
		matches := paramRe.FindAllStringIndex(template.RawSQL, -1)

		var regexStr strings.Builder
		regexStr.WriteString("(?i)")

		lastIndex := 0

		for _, match := range matches {
			start, end := match[0], match[1]
			regexStr.WriteString(regexp.QuoteMeta(template.RawSQL[lastIndex:start]))
			regexStr.WriteString(`(.+?)`)
			lastIndex = end
		}
		regexStr.WriteString(regexp.QuoteMeta(template.RawSQL[lastIndex:]))
		regexStr.WriteString(`\s*;?\s*$`)

		compiledRe, err := regexp.Compile("^" + regexStr.String())
		if err != nil {
			return nil, err
		}

		r.mu.Lock()
		r.regexCache[template.RawSQL] = compiledRe
		r.mu.Unlock()
		re = compiledRe
	}

	matches := re.FindStringSubmatch(sql)
	if matches == nil {
		return nil, fmt.Errorf("SQL does not match template regex")
	}

	paramRe := regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)
	paramMatches := paramRe.FindAllString(template.RawSQL, -1)

	if len(matches)-1 != len(paramMatches) {
		return nil, fmt.Errorf("match count mismatch")
	}

	result := make(map[string]interface{})
	for i, pName := range paramMatches {
		val := matches[i+1]
		val = strings.Trim(val, "'")
		result[pName] = val
	}

	return result, nil
}

// normalizeQuery provides a basic normalization for lookup.
func normalizeQuery(q string) string {
	q = strings.ToLower(q)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
	q = regexp.MustCompile(`'[^']*'`).ReplaceAllString(q, "?")
	q = regexp.MustCompile(`\d+`).ReplaceAllString(q, "?")
	q = regexp.MustCompile(`:\w+`).ReplaceAllString(q, "?")
	q = regexp.MustCompile(`=\s*\?`).ReplaceAllString(q, " = ?")
	return strings.TrimSpace(q)
}
