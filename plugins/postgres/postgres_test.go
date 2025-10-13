package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/database"
)

func TestPostgresPlugin(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_MACHINE_NAME") == "" {
		t.Skip("Skipping test; Docker not available")
	}

	dsn := "postgres://user:password@localhost:5432/test?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, db.PingContext(ctx))

	// Test schema extraction
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test_table (id INT PRIMARY KEY, name VARCHAR(255))")
	require.NoError(t, err)

	extractor := database.NewPostgresSchemaExtractor(db)
	schema, err := extractor.ExtractSchema(ctx)
	require.NoError(t, err)
	assert.Equal(t, "test", schema.Name)
	assert.Contains(t, schema.Tables, "test_table")

	// Test schema loading
	loader := database.NewPostgresSchemaLoader(db)
	loadSchema := &models.DatabaseSchema{
		Name: "test",
		Tables: map[string]*models.TableSchema{
			"load_test_table": {
				Name: "load_test_table",
				Columns: []*models.ColumnSchema{
					{Name: "id", Type: "integer"},
					{Name: "value", Type: "text"},
				},
				PK: []string{"id"},
			},
		},
	}
	err = loader.LoadSchema(ctx, loadSchema)
	require.NoError(t, err)

	rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_name = 'load_test_table'")
	require.NoError(t, err)
	defer rows.Close()
	assert.True(t, rows.Next())
}