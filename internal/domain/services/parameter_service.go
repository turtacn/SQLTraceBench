package services

import (
	"fmt"
	"regexp"
	"sort"
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
}

// NewParameterService creates a new ParameterService.
func NewParameterService() *ParameterService {
	return &ParameterService{
		extractor: &RegexParameterExtractor{
			regexCache: make(map[string]*regexp.Regexp),
		},
	}
}

// BuildModel analyzes a collection of SQL traces and builds a ParameterModel.
func (s *ParameterService) BuildModel(tc models.TraceCollection, templates []models.SQLTemplate) *models.ParameterModel {
	pm := models.NewParameterModel()

	// If there are no traces, create a default model from the templates.
	if len(tc.Traces) == 0 {
		for _, t := range templates {
			// Ensure parameters are extracted/known
			if len(t.Parameters) == 0 {
				t.ExtractParameters()
			}

			if len(t.Parameters) > 0 {
				if _, ok := pm.TemplateParameters[t.GroupKey]; !ok {
					pm.TemplateParameters[t.GroupKey] = make(map[string]*models.ValueDistribution)
				}
				for _, pName := range t.Parameters {
					// Create a default distribution for each parameter.
					dist := models.NewValueDistribution()
					// For default model, maybe we don't add dummy observations,
					// or we add a placeholder.
					dist.AddObservation("1")
					pm.TemplateParameters[t.GroupKey][pName] = dist
				}
			}
		}
		return pm
	}

	templateMap := make(map[string]models.SQLTemplate)
	for _, t := range templates {
		// Ensure parameters are extracted
		if len(t.Parameters) == 0 {
			t.ExtractParameters() // This modifies the copy in the loop if not careful?
			// No, range over slice returns copy. But we might need the parameters later.
			// Actually t is a copy. But ExtractParameters modifies fields of the struct.
			// Since we put it in map, the one in map will have it.
		}
		// Better to ensure the original slice has them or we use pointers.
		// models.SQLTemplate is a struct.

		// Let's explicitly fix the template parameters on the copy we store.
		temp := t
		temp.ExtractParameters()
		templateMap[normalizeQuery(temp.RawSQL)] = temp
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

// RegexParameterExtractor extracts parameter values using regex matching against the template.
type RegexParameterExtractor struct {
	regexCache map[string]*regexp.Regexp
	mu         sync.RWMutex
}

func (r *RegexParameterExtractor) Extract(sql string, template *models.SQLTemplate) (map[string]interface{}, error) {
	// We need to convert the template RawSQL into a regex that captures values.
	// Template: SELECT * FROM users WHERE id = :id AND name = :name
	// Regex:    SELECT * FROM users WHERE id = (.*?) AND name = (.*?)

	// Check cache
	r.mu.RLock()
	re, ok := r.regexCache[template.RawSQL]
	r.mu.RUnlock()

	if !ok {
		// Build regex
		// Escape the template SQL to be safe for regex, except for the parameters
		// This is tricky. simpler approach:
		// 1. Identify parameter positions.
		// 2. Replace parameters with capture groups.
		// 3. Escape everything else.

		// To do this reliably, we need to split by parameters.
		// But parameters are identified by :name.

		// Let's simplify: replace all :paramName with (.*?) or similar.
		// But first escape special regex characters in the static parts.

		// pattern := regexp.QuoteMeta(template.RawSQL) // Unused

		// Now we have escaped version. We need to unescape the :paramName parts and replace them with capture groups.
		// But wait, QuoteMeta escapes ':' too? No, ':' is not special in regex usually, but let's check.
		// Go's QuoteMeta escapes: \.+*?()|[]{}^$

		// So :paramName remains :paramName.
		// We can iterate over parameters (sorted by length descending to avoid partial matches?)
		// and replace them with capture group.

		// Sort parameters by length descending to avoid replacing ":prefix" inside ":prefix_suffix"
		sortedParams := make([]string, len(template.Parameters))
		copy(sortedParams, template.Parameters)
		sort.Slice(sortedParams, func(i, j int) bool {
			return len(sortedParams[i]) > len(sortedParams[j])
		})

		// We also need to map capture group index to parameter name.
		// The order of capture groups depends on where they appear in the string.
		// So we can't just iterate parameters and replace.
		// We need to find their positions.

		// Alternative: Construct the regex from scratch.
		// Using the paramRe from models (which isn't exported, but we can recreate or assume :w+)

		// Let's try a simple replace.
		// Note: This naive regex approach fails if parameters are repeated or order is complex,
		// but it's better than the mock.

		// Map to store the order of parameters as they appear in the regex
		// We can't easily know the order unless we parse.
		// But we can assume the extraction gives us submatches in order.

		// Let's use a placeholder implementation that assumes simple cases for now,
		// or try to find all params and replace them.

		// Robust approach:
		// Find all parameter occurrences in the template string with their indices.
		// Sort them by index.
		// Build the regex string by appending parts.

		// Re-using the regex from models/template.go logic
		paramRe := regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)
		matches := paramRe.FindAllStringIndex(template.RawSQL, -1)

		var regexStr strings.Builder
		regexStr.WriteString("(?i)") // Case insensitive matching for SQL keywords

		lastIndex := 0
		paramOrder := make([]string, 0)

		for _, match := range matches {
			start, end := match[0], match[1]
			// Append the static part before this parameter, escaped
			regexStr.WriteString(regexp.QuoteMeta(template.RawSQL[lastIndex:start]))

			// Append capture group
			// We use '.*?' to capture lazily until the next part.
			// Depending on SQL, value might be quoted or a number.
			// (\d+|'[^']*'|[^ ,]+) might be better?
			// Let's stick to (.*?) but it might match too much or too little.
			// Ideally: (\S+) or similar.
			// Let's use `(.+?)`
			regexStr.WriteString(`(.+?)`)

			paramName := template.RawSQL[start:end]
			paramOrder = append(paramOrder, paramName)

			lastIndex = end
		}
		regexStr.WriteString(regexp.QuoteMeta(template.RawSQL[lastIndex:]))
		// Allow some flexibility with trailing spaces/semicolons
		regexStr.WriteString(`\s*;?\s*$`)

		compiledRe, err := regexp.Compile("^" + regexStr.String())
		if err != nil {
			return nil, err
		}

		// We need to store paramOrder with the regex to map back.
		// For now, we can re-derive it or assume it's consistent.
		// But regexCache just stores *regexp.Regexp.
		// We need a struct or just re-derive.
		// Re-deriving is fast enough.

		r.mu.Lock()
		r.regexCache[template.RawSQL] = compiledRe
		r.mu.Unlock()
		re = compiledRe
	}

	// Now match
	matches := re.FindStringSubmatch(sql)
	if matches == nil {
		return nil, fmt.Errorf("SQL does not match template regex")
	}

	// Re-determine parameter order to map values
	// (This is duplicated effort, ideally cache this too)
	paramRe := regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)
	paramMatches := paramRe.FindAllString(template.RawSQL, -1)

	if len(matches)-1 != len(paramMatches) {
		// This might happen if a parameter appears multiple times?
		// Or if our regex construction logic was flawed.
		return nil, fmt.Errorf("match count mismatch")
	}

	result := make(map[string]interface{})
	for i, pName := range paramMatches {
		val := matches[i+1]
		// Clean up quotes if necessary?
		// For now, keep as string.
		val = strings.Trim(val, "'") // simple cleanup
		result[pName] = val
	}

	return result, nil
}

// normalizeQuery provides a basic normalization for lookup.
func normalizeQuery(q string) string {
	q = strings.ToLower(q)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
	// Normalize values to ?
	q = regexp.MustCompile(`'[^']*'`).ReplaceAllString(q, "?")
	q = regexp.MustCompile(`\d+`).ReplaceAllString(q, "?")
	// Normalize params to ?
	q = regexp.MustCompile(`:\w+`).ReplaceAllString(q, "?")
	// Normalize = ?
	q = regexp.MustCompile(`=\s*\?`).ReplaceAllString(q, " = ?")

	return strings.TrimSpace(q)
}
