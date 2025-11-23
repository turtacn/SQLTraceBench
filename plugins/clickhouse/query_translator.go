package clickhouse

import (
	"regexp"
	"strings"
)

// QueryTranslator is the interface for translating SQL queries.
type QueryTranslator interface {
	TranslateQuery(sql string) (string, error)
}

type ClickHouseTranslator struct{}

// NewQueryTranslator creates a new QueryTranslator.
func NewQueryTranslator() QueryTranslator {
	return &ClickHouseTranslator{}
}

// TranslateQuery translates a generic or MySQL-specific SQL query to ClickHouse dialect.
func (t *ClickHouseTranslator) TranslateQuery(sql string) (string, error) {
	// 1. Sanitize backticks (MySQL) to double quotes or nothing?
	// ClickHouse uses backticks for identifiers too, so actually MySQL backticks are usually fine in ClickHouse!
	// But the prompt says: "Replace backticks ` for ClickHouse legal symbols".
	// ClickHouse documentation says: "Identifiers are quoted with backticks or double quotes."
	// Maybe the prompt implies removing them or handling specific cases.
	// "TestSanitize: Input `SELECT * FROM `User`` (MySQL style), assert output fits ClickHouse (clean backticks)."
	// If the table name is `User` (case sensitive?), ClickHouse is case sensitive.
	// Let's assume the requirement is to keep them or convert to double quotes if needed.
	// However, prompt example says "Replace backticks ` for ClickHouse legal symbols".
	// Wait, standard SQL uses double quotes for identifiers. MySQL uses backticks.
	// ClickHouse supports both.
	// But if the prompt specifically asks to "clean" them, maybe it means removing them for simple identifiers?
	// Or maybe converting to double quotes. I'll stick to replacing backticks with nothing if simple, or double quotes.
	// Let's look at `TestSanitize` expectation in prompt: "TestSanitize: Input `SELECT * FROM `User`` (MySQL style), assert output fits ClickHouse (cleaning backticks)."
	// "cleaning backticks" sounds like removing them.

	// 2. Remove ENGINE=InnoDB
	// Regex to remove ENGINE=...

	translated := sql

	// Remove ENGINE=...
	reEngine := regexp.MustCompile(`(?i)\s*ENGINE\s*=\s*\w+`)
	translated = reEngine.ReplaceAllString(translated, "")

	// Remove CHARSET=...
	reCharset := regexp.MustCompile(`(?i)\s*DEFAULT\s*CHARSET\s*=\s*\w+`)
	translated = reCharset.ReplaceAllString(translated, "")

	reCharsetSimple := regexp.MustCompile(`(?i)\s*CHARSET\s*=\s*\w+`)
	translated = reCharsetSimple.ReplaceAllString(translated, "")

	// Replace backticks with double quotes (standard SQL, supported by CH)
	// Or remove them if they are just around simple identifiers.
	// Replacing with double quotes is safer.
	translated = strings.ReplaceAll(translated, "`", "\"")

	// ClickHouse specific: AUTO_INCREMENT is not supported in the same way (usually)
	reAutoInc := regexp.MustCompile(`(?i)\s*AUTO_INCREMENT`)
	translated = reAutoInc.ReplaceAllString(translated, "")

	// Clean up multiple spaces
	reSpaces := regexp.MustCompile(`\s+`)
	translated = reSpaces.ReplaceAllString(translated, " ")

	return strings.TrimSpace(translated), nil
}
