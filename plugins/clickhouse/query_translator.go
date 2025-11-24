package clickhouse

import (
	"strings"
)

// QueryTranslator is the interface for query translation logic.
type QueryTranslator interface {
	TranslateQuery(sql string) (string, error)
}

type ClickHouseTranslator struct{}

// NewQueryTranslator creates a new QueryTranslator.
func NewQueryTranslator() QueryTranslator {
	return &ClickHouseTranslator{}
}

// TranslateQuery translates a MySQL query to a ClickHouse query.
func (t *ClickHouseTranslator) TranslateQuery(sql string) (string, error) {
	processed := sql

	// 1. Remove SQL_NO_CACHE hint
	processed = strings.ReplaceAll(processed, "SQL_NO_CACHE", "")

	// 2. Keep backticks (ClickHouse supports them, so no action needed)

	// 3. Normalize function names: NOW() -> now()
	// Using simple replace for now as per requirement MVP.
	// Case insensitive replacement would be better but strings.ReplaceAll is case sensitive.
	// We can use regex if we want robust case handling, but "NOW()" usually appears in uppercase in generated queries.
	// Let's iterate over a few common functions if needed, but requirement example is explicit: NOW() -> now()
	processed = strings.ReplaceAll(processed, "NOW()", "now()")

	// 4. Remove trailing semicolon
	processed = strings.TrimSpace(processed)
	if strings.HasSuffix(processed, ";") {
		processed = processed[:len(processed)-1]
	}

	// 4. Remove extra spaces
	processed = strings.Join(strings.Fields(processed), " ")

	return processed, nil
}
