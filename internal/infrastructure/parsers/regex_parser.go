// Package parsers contains the concrete implementations of the Parser interface.
package parsers

import (
	"regexp"
	"strings"
)

// re is the regular expression used to find table names in SQL queries.
// It specifically looks for tables following `FROM` or `JOIN` clauses.
var re = regexp.MustCompile(`(?is)(?:from|join)\s+([a-zA-Z0-9_\.]+)`)

// RegexParser is a SQL parser that uses regular expressions to extract information from queries.
type RegexParser struct{}

// NewRegexParser creates a new RegexParser.
func NewRegexParser() *RegexParser {
	return &RegexParser{}
}

// ListTables extracts table names from a SQL query using a regular expression.
// It finds all matches for tables in `FROM` and `JOIN` clauses and returns a deduplicated slice of table names.
func (p *RegexParser) ListTables(sql string) ([]string, error) {
	matches := re.FindAllStringSubmatch(sql, -1)
	if matches == nil {
		return []string{}, nil
	}

	seen := make(map[string]struct{})
	tables := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			tableName := strings.TrimSpace(match[1])
			if _, ok := seen[tableName]; !ok {
				seen[tableName] = struct{}{}
				tables = append(tables, tableName)
			}
		}
	}
	return tables, nil
}