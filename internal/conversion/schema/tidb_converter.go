package schema

import (
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// TiDBConverter implements SchemaConverter for TiDB.
// Since TiDB is MySQL compatible, it reuses MySQLConverter logic but handles TiDB specifics.
type TiDBConverter struct {
	mysqlConverter *MySQLConverter
}

// NewTiDBConverter creates a new TiDBConverter.
func NewTiDBConverter() *TiDBConverter {
	return &TiDBConverter{
		mysqlConverter: NewMySQLConverter(),
	}
}

// ConvertDDL converts TiDB DDL to target DB format.
func (c *TiDBConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
	// Pre-process TiDB specific syntax

	// 1. Remove /*T![clustered_index] ... */
	// 2. Remove SHARD_ROW_ID_BITS=N
	// 3. Remove AUTO_RANDOM

	cleanDDL := c.preprocessTiDBDDL(sourceDDL)

	return c.mysqlConverter.ConvertDDL(cleanDDL, targetDB)
}

func (c *TiDBConverter) preprocessTiDBDDL(ddl string) string {
	// Remove SHARD_ROW_ID_BITS
	// Replace AUTO_RANDOM with nothing (or convert to normal int if needed, but AUTO_INCREMENT removal handles it)
	// But AUTO_RANDOM is an attribute on column.

	lines := strings.Split(ddl, "\n")
	var sb strings.Builder
	for _, line := range lines {
		// Simple line based removal for table options
		if strings.Contains(line, "SHARD_ROW_ID_BITS") {
			// Check if statement ends with semicolon
			hasSemicolon := strings.TrimSpace(line)[len(strings.TrimSpace(line))-1] == ';'

			idx := strings.Index(line, "SHARD_ROW_ID_BITS")
			line = line[:idx]

			if hasSemicolon {
				line = strings.TrimRight(line, " \t") + ";"
			}
		}

		// Remove CLUSTERED INDEX comments
		// /*T![clustered_index] CLUSTERED */
		line = strings.ReplaceAll(line, "/*T![clustered_index] CLUSTERED */", "")
		line = strings.ReplaceAll(line, "/*T![clustered_index] NONCLUSTERED */", "")

		// Remove AUTO_RANDOM
		line = strings.ReplaceAll(line, "AUTO_RANDOM", "")

		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

// ConvertTable converts a single table structure.
func (c *TiDBConverter) ConvertTable(sourceTable *models.TableSchema, targetDB string) (*models.TableSchema, error) {
	// Delegate to MySQL converter
	return c.mysqlConverter.ConvertTable(sourceTable, targetDB)
}

// GetTypeMapping gets the type mapping.
func (c *TiDBConverter) GetTypeMapping(sourceType string, targetDB string) (string, error) {
	return c.mysqlConverter.GetTypeMapping(sourceType, targetDB)
}
