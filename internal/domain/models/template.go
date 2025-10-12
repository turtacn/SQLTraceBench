package models

import (
	"fmt"
	"regexp"
	"sort"
)

// paramRe is a regular expression used to identify named parameters in a SQL query.
var paramRe = regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)

// SQLTemplate represents a normalized SQL query with its parameters extracted.
type SQLTemplate struct {
	RawSQL     string
	GroupKey   string
	Weight     int
	Parameters []string
}

// ExtractParameters finds all named parameters in the RawSQL query.
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

// GenerateQuery creates a QueryWithArgs struct from the template.
// It replaces named parameters with '?' placeholders and populates the args slice.
func (t *SQLTemplate) GenerateQuery(params map[string]interface{}) (QueryWithArgs, error) {
	query := t.RawSQL
	args := make([]interface{}, 0, len(t.Parameters))

	// This is a simplified approach that assumes parameter order.
	// This is a known limitation to be addressed with a better parser.
	for _, pName := range t.Parameters {
		val, ok := params[pName]
		if !ok {
			return QueryWithArgs{}, fmt.Errorf("parameter %s not found in params map", pName)
		}
		args = append(args, val)
		query = regexp.MustCompile(pName).ReplaceAllString(query, "?")
	}

	return QueryWithArgs{Query: query, Args: args}, nil
}