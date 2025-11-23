package starrocks

import (
	"regexp"
	"strings"
)

// StarRocksTranslator handles query translation from MySQL dialect to StarRocks dialect.
type StarRocksTranslator struct{}

// Translate adapts the SQL query for StarRocks.
func (t *StarRocksTranslator) Translate(sql string) (string, error) {
	// Remove SQL_NO_CACHE
	reNoCache := regexp.MustCompile(`(?i)\bSQL_NO_CACHE\b`)
	sql = reNoCache.ReplaceAllString(sql, "")

	// Remove SQL_CALC_FOUND_ROWS
	reCalc := regexp.MustCompile(`(?i)\bSQL_CALC_FOUND_ROWS\b`)
	sql = reCalc.ReplaceAllString(sql, "")

	// Remove FORCE INDEX hints as they might not be compatible or needed
	// Example: FORCE INDEX (idx_name)
	reForceIndex := regexp.MustCompile(`(?i)\bFORCE\s+INDEX\s*\([^\)]+\)`)
	sql = reForceIndex.ReplaceAllString(sql, "")

	// Clean up extra whitespace
	reSpaces := regexp.MustCompile(`\s+`)
	sql = reSpaces.ReplaceAllString(sql, " ")

	return strings.TrimSpace(sql), nil
}
