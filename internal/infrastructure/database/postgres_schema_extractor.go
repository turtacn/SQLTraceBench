package database

import (
	"context"
	"database/sql"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/types"
)

// PostgresSchemaExtractor is a schema extractor for PostgreSQL databases.
type PostgresSchemaExtractor struct {
	db *sql.DB
}

// NewPostgresSchemaExtractor creates a new PostgresSchemaExtractor.
func NewPostgresSchemaExtractor(db *sql.DB) *PostgresSchemaExtractor {
	return &PostgresSchemaExtractor{db: db}
}

// ExtractSchema extracts the schema from a PostgreSQL database.
func (e *PostgresSchemaExtractor) ExtractSchema(ctx context.Context) (*models.DatabaseSchema, error) {
	var dbName string
	if err := e.db.QueryRowContext(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get database name", err)
	}

	schema := &models.DatabaseSchema{
		Name:   dbName,
		Tables: make([]*models.TableSchema, 0),
	}

	rows, err := e.db.QueryContext(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`)
	if err != nil {
		return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get tables", err)
	}
	defer rows.Close()

	var tables []*models.TableSchema
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, types.WrapError(types.ErrDatabaseConnection, "failed to scan table name", err)
		}
		tables = append(tables, &models.TableSchema{Name: tableName})
	}

	for _, table := range tables {
		colRows, err := e.db.QueryContext(ctx, `
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = $1
		`, table.Name)
		if err != nil {
			return nil, types.WrapError(types.ErrDatabaseConnection, "failed to get columns", err)
		}
		defer colRows.Close()

		for colRows.Next() {
			var col models.ColumnSchema
			var isNullable string
			var colDefault sql.NullString
			if err := colRows.Scan(&col.Name, &col.DataType, &isNullable, &colDefault); err != nil {
				return nil, types.WrapError(types.ErrDatabaseConnection, "failed to scan column", err)
			}
			col.IsNullable = isNullable == "YES"
			col.Default = colDefault.String
			table.Columns = append(table.Columns, &col)
		}
	}
	schema.Tables = tables

	return schema, nil
}