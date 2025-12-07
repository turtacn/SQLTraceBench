package schema

import (
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// TiDBConverter implements SchemaConverter for TiDB.
// Since TiDB is MySQL compatible, it reuses MySQLConverter logic but handles TiDB specifics.
type TiDBConverter struct {
	mysqlConverter *MySQLConverter
    // TiDB converter could have its own type mapper if needed, but reusing MySQL's is fine for now
    // as TiDB types are mostly MySQL compatible.
}

// NewTiDBConverter creates a new TiDBConverter.
func NewTiDBConverter() *TiDBConverter {
	return &TiDBConverter{
		mysqlConverter: NewMySQLConverter(),
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
    // For intelligent mapping, we might want to override the SourceDB to "tidb"
    // in the context. Since MySQLConverter hardcodes "mysql", we might need to modify MySQLConverter
    // or manually handle it here.

    // If we want full "tidb" source support in rules, we should copy logic from MySQLConverter but use "tidb".
    // Or we can modify MySQLConverter to accept sourceDB name.

    // For now, reusing MySQL converter is acceptable as per requirements ("Integration" task lists it).
    // The previous implementation simply delegated.

	return c.mysqlConverter.ConvertTable(sourceTable, targetDB)
}

// GetTypeMapping gets the type mapping.
func (c *TiDBConverter) GetTypeMapping(sourceType string, targetDB string) (string, error) {
	return c.mysqlConverter.GetTypeMapping(sourceType, targetDB)
}
