package models

import (
	"fmt"
	"regexp"
	"sort"
)

// paramRe is a regular expression used to identify named parameters in a SQL query.
// It looks for patterns like `:param_name`.
var paramRe = regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)

// SQLTemplate represents a normalized SQL query with its parameters extracted.
// It serves as a blueprint for generating new queries with different parameter values.
type SQLTemplate struct {
	// RawSQL is the original, un-normalized SQL query.
	RawSQL string
	// GroupKey is a normalized version of the SQL query, used for grouping similar queries.
	GroupKey string
	// Weight is the frequency of this template in the trace, used for workload generation.
	Weight int
	// Parameters is a list of named parameters found in the query.
	Parameters []string
}

// ExtractParameters finds all named parameters in the RawSQL query and populates the Parameters field.
// The extracted parameters are sorted to ensure a consistent order.
func (t *SQLTemplate) ExtractParameters() {
	params := paramRe.FindAllString(t.RawSQL, -1)
	paramSet := make(map[string]struct{})
	for _, p := range params {
		paramSet[p] = struct{}{}
	}

	t.Parameters = make([]string, 0, len(paramSet))
	for p := range paramSet {
		t.Parameters = append(t.Parameters, p)
	}
	sort.Strings(t.Parameters)
}

// GenerateQuery creates a new SQL query from the template by substituting the named parameters
// with the provided values. It performs basic type checking to ensure that string values are
// quoted correctly.
func (t *SQLTemplate) GenerateQuery(params map[string]interface{}) (string, error) {
	query := t.RawSQL
	for _, p := range t.Parameters {
		val, ok := params[p]
		if !ok {
			return "", fmt.Errorf("parameter %s not found in params map", p)
		}

		var replacement string
		switch v := val.(type) {
		case string:
			replacement = fmt.Sprintf("'%s'", v) // Basic string quoting
		case int, int64, float64:
			replacement = fmt.Sprintf("%v", v)
		default:
			return "", fmt.Errorf("unsupported parameter type for %s: %T", p, v)
		}
		query = regexp.MustCompile(p).ReplaceAllString(query, replacement)
	}
	return query, nil
}