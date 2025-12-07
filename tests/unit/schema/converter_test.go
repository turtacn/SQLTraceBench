package schema_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/conversion/schema"
)

func normalizeSQL(sql string) string {
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	for strings.Contains(sql, "  ") {
		sql = strings.ReplaceAll(sql, "  ", " ")
	}
	return strings.TrimSpace(sql)
}

func TestMySQLConverter_BasicTypes(t *testing.T) {
	converter := schema.NewMySQLConverter()
	testCases := []struct {
		name      string
		inputDDL  string
		targetDB  string
		expectDDL string
		expectErr bool
	}{
		{
			name:      "Simple INT to ClickHouse",
			inputDDL:  "CREATE TABLE users (id INT, name VARCHAR(100));",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE users ( id Int32, name String ) ENGINE = MergeTree() ORDER BY tuple();",
		},
		{
			name:      "PrimaryKey to OrderBy",
			inputDDL:  "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE users ( id Int32, name String ) ENGINE = MergeTree() ORDER BY (id);",
		},
		{
			name:      "ENUM type conversion",
			inputDDL:  "CREATE TABLE orders (status ENUM('pending','shipped'));",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE orders ( status Enum8('pending'=1, 'shipped'=2) ) ENGINE = MergeTree() ORDER BY tuple();",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.ConvertDDL(tc.inputDDL, tc.targetDB)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, normalizeSQL(tc.expectDDL), normalizeSQL(result))
			}
		})
	}
}

func TestPostgresConverter_BasicTypes(t *testing.T) {
	converter := schema.NewPostgresConverter()
	testCases := []struct {
		name      string
		inputDDL  string
		targetDB  string
		expectDDL string
		expectErr bool
	}{
		{
			name:      "Simple INT to ClickHouse",
			inputDDL:  "CREATE TABLE users (id INTEGER, name VARCHAR(100));",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE users ( id Int32, name String ) ENGINE = MergeTree() ORDER BY tuple();",
		},
		{
			name:      "Serial PK",
			inputDDL:  "CREATE TABLE items (id SERIAL PRIMARY KEY, info TEXT);",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE items ( id Int32, info String ) ENGINE = MergeTree() ORDER BY (id);",
		},
		{
			name:      "Array Type",
			inputDDL:  "CREATE TABLE posts (tags INTEGER[]);",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE posts ( tags Array(Int32) ) ENGINE = MergeTree() ORDER BY tuple();",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.ConvertDDL(tc.inputDDL, tc.targetDB)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, normalizeSQL(tc.expectDDL), normalizeSQL(result))
			}
		})
	}
}

func TestTiDBConverter_Specifics(t *testing.T) {
	converter := schema.NewTiDBConverter()
	testCases := []struct {
		name      string
		inputDDL  string
		targetDB  string
		expectDDL string
		expectErr bool
	}{
		{
			name:      "ShardRowID Removal",
			inputDDL:  "CREATE TABLE t (a INT) SHARD_ROW_ID_BITS=4;",
			targetDB:  "clickhouse",
			expectDDL: "CREATE TABLE t ( a Int32 ) ENGINE = MergeTree() ORDER BY tuple();",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.ConvertDDL(tc.inputDDL, tc.targetDB)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, normalizeSQL(tc.expectDDL), normalizeSQL(result))
			}
		})
	}
}
