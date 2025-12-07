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
		mysqlConverter: NewMySQLConverterWithSource("tidb"),
	}
}

// ConvertDDL converts TiDB DDL to target DB format.
func (c *TiDBConverter) ConvertDDL(sourceDDL string, targetDB string) (string, error) {
	cleanDDL := c.preprocessTiDBDDL(sourceDDL)
	return c.mysqlConverter.ConvertDDL(cleanDDL, targetDB)
}

func (c *TiDBConverter) preprocessTiDBDDL(ddl string) string {
	lines := strings.Split(ddl, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if strings.Contains(line, "SHARD_ROW_ID_BITS") {
			hasSemicolon := strings.TrimSpace(line)[len(strings.TrimSpace(line))-1] == ';'
			idx := strings.Index(line, "SHARD_ROW_ID_BITS")
			line = line[:idx]

			if hasSemicolon {
				line = strings.TrimRight(line, " \t,") + ";"
			} else {
                 line = strings.TrimRight(line, " \t,")
            }
		}

		line = strings.ReplaceAll(line, "/*T![clustered_index] CLUSTERED */", "")
		line = strings.ReplaceAll(line, "/*T![clustered_index] NONCLUSTERED */", "")
		line = strings.ReplaceAll(line, "AUTO_RANDOM", "")

		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

// ConvertTable converts a single table structure.
func (c *TiDBConverter) ConvertTable(sourceTable *models.TableSchema, targetDB string) (*models.TableSchema, error) {
	return c.mysqlConverter.ConvertTable(sourceTable, targetDB)
}

// GetTypeMapping gets the type mapping.
func (c *TiDBConverter) GetTypeMapping(sourceType string, targetDB string) (string, error) {
	return c.mysqlConverter.GetTypeMapping(sourceType, targetDB)
}
