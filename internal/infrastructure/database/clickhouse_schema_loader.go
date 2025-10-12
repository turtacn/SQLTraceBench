package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// ClickHouseSchemaLoader is a schema loader for ClickHouse databases.
type ClickHouseSchemaLoader struct {
	db *sql.DB
}

// NewClickHouseSchemaLoader creates a new ClickHouseSchemaLoader.
func NewClickHouseSchemaLoader(db *sql.DB) *ClickHouseSchemaLoader {
	return &ClickHouseSchemaLoader{db: db}
}

// LoadSchema loads a schema into a ClickHouse database.
func (l *ClickHouseSchemaLoader) LoadSchema(ctx context.Context, schema *models.DatabaseSchema) error {
	for _, table := range schema.Tables {
		ddl, err := l.generateCreateTableDDL(table)
		if err != nil {
			return err
		}

		if _, err := l.db.ExecContext(ctx, ddl); err != nil {
			return types.WrapError(types.ErrDatabaseConnection, fmt.Sprintf("failed to create table %s", table.Name), err)
		}
	}
	return nil
}

func (l *ClickHouseSchemaLoader) generateCreateTableDDL(table *models.TableSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", table.Name))

	for i, col := range table.Columns {
		sb.WriteString(fmt.Sprintf("  `%s` %s", col.Name, col.Type))
		if i < len(table.Columns)-1 {
			sb.WriteString(",\n")
		}
	}

	sb.WriteString(fmt.Sprintf("\n) ENGINE = MergeTree() ORDER BY (%s);", strings.Join(table.PK, ",")))
	return sb.String(), nil
}