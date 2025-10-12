package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// MySqlSchemaLoader is a schema loader for MySQL databases.
type MySqlSchemaLoader struct {
	db *sql.DB
}

// NewMySqlSchemaLoader creates a new MySqlSchemaLoader.
func NewMySqlSchemaLoader(db *sql.DB) *MySqlSchemaLoader {
	return &MySqlSchemaLoader{db: db}
}

// LoadSchema loads a schema into a MySQL database.
func (l *MySqlSchemaLoader) LoadSchema(ctx context.Context, schema *models.DatabaseSchema) error {
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

func (l *MySqlSchemaLoader) generateCreateTableDDL(table *models.TableSchema) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", table.Name))

	for i, col := range table.Columns {
		sb.WriteString(fmt.Sprintf("  `%s` %s", col.Name, col.Type))
		if !col.IsNullable {
			sb.WriteString(" NOT NULL")
		}
		if col.Default != "" {
			sb.WriteString(fmt.Sprintf(" DEFAULT '%s'", col.Default))
		}
		if i < len(table.Columns)-1 || len(table.PK) > 0 {
			sb.WriteString(",\n")
		}
	}

	if len(table.PK) > 0 {
		sb.WriteString(fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(table.PK, ",")))
	}

	sb.WriteString("\n);")
	return sb.String(), nil
}