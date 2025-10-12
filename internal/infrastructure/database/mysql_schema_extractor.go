package database

import (
	"context"
	"database/sql"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// MySqlSchemaExtractor is a schema extractor for MySQL databases.
type MySqlSchemaExtractor struct {
	db *sql.DB
}

// NewMySqlSchemaExtractor creates a new MySqlSchemaExtractor.
func NewMySqlSchemaExtractor(db *sql.DB) *MySqlSchemaExtractor {
	return &MySqlSchemaExtractor{db: db}
}

// ExtractSchema extracts the schema from a MySQL database.
func (e *MySqlSchemaExtractor) ExtractSchema(ctx context.Context) (*models.DatabaseSchema, error) {
	// Get the database name.
	var dbName string
	if err := e.db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName); err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get database name", err)
	}

	schema := &models.DatabaseSchema{
		Name:   dbName,
		Tables: make(map[string]*models.TableSchema),
	}

	// Get the tables.
	rows, err := e.db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = ?", dbName)
	if err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get tables", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, types.WrapError(types.ErrDatabaseConnection, "failed to scan table name", err)
		}
		schema.Tables[tableName] = &models.TableSchema{Name: tableName}
	}

	// Get the columns for each table.
	for tableName, table := range schema.Tables {
		colRows, err := e.db.QueryContext(ctx, `
			SELECT column_name, column_type, is_nullable, column_default
			FROM information_schema.columns
			WHERE table_schema = ? AND table_name = ?
		`, dbName, tableName)
		if err != nil {
			return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get columns", err)
		}
		defer colRows.Close()

		for colRows.Next() {
			var col models.ColumnSchema
			var isNullable string
			var colDefault sql.NullString
			if err := colRows.Scan(&col.Name, &col.Type, &isNullable, &colDefault); err != nil {
				return nil, types.WrapError(types.ErrDatabaseConnection, "failed to scan column", err)
			}
			col.IsNullable = isNullable == "YES"
			col.Default = colDefault.String
			table.Columns = append(table.Columns, &col)
		}
	}

	return schema, nil
}