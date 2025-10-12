package cmd

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaMigration(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_MACHINE_NAME") == "" {
		t.Skip("Skipping test; Docker not available")
	}

	mysqlDsn := "root:root@tcp(127.0.0.1:3306)/test?parseTime=true"
	mysqlDb, err := sql.Open("mysql", mysqlDsn)
	require.NoError(t, err)
	defer mysqlDb.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, mysqlDb.PingContext(ctx))

	_, err = mysqlDb.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255))")
	require.NoError(t, err)

	tmpDir, err := ioutil.TempDir("", "schema_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	schemaPath := filepath.Join(tmpDir, "schema.json")
	rootCmd.SetArgs([]string{"schema", "dump", "--dsn", mysqlDsn, "--out", schemaPath})
	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.FileExists(t, schemaPath)

	clickhouseDsn := "tcp://127.0.0.1:9000"
	rootCmd.SetArgs([]string{
		"schema", "load",
		"--schema", schemaPath,
		"--dsn", clickhouseDsn,
		"--driver", "clickhouse",
		"--target", "clickhouse",
	})
	err = rootCmd.Execute()
	require.NoError(t, err)

	// Verify the schema in ClickHouse
	chDb, err := sql.Open("clickhouse", clickhouseDsn)
	require.NoError(t, err)
	defer chDb.Close()
	require.NoError(t, chDb.PingContext(ctx))

	rows, err := chDb.QueryContext(ctx, "DESCRIBE TABLE users")
	require.NoError(t, err)
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var name, type_, defaultType, defaultExpr, comment, codec, ttl string
		err := rows.Scan(&name, &type_, &defaultType, &defaultExpr, &comment, &codec, &ttl)
		require.NoError(t, err)
		columns = append(columns, name)
	}
	assert.ElementsMatch(t, []string{"id", "name", "email"}, columns)
}