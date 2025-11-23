package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestConvertType(t *testing.T) {
	c := NewSchemaConverter().(*ClickHouseConverter)

	tests := []struct {
		input    string
		expected string
	}{
		{"tinyint", "Int8"},
		{"int", "Int32"},
		{"bigint", "Int64"},
		{"varchar(255)", "String"},
		{"DECIMAL(10,2)", "Decimal64(2)"},
		{"datetime", "DateTime64"},
		{"timestamp", "DateTime64"},
		{"date", "Date32"},
		{"boolean", "Bool"},
		{"unknown", "String"},
	}

	for _, test := range tests {
		result := c.mapType(test.input)
		assert.Equal(t, test.expected, result, "Failed to map type: %s", test.input)
	}
}

func TestConvertSchema_PrimaryKey(t *testing.T) {
	c := NewSchemaConverter().(*ClickHouseConverter)

	srcSchema := &models.Schema{
		Databases: []models.DatabaseSchema{
			{
				Name: "test_db",
				Tables: []*models.TableSchema{
					{
						Name: "users",
						Columns: []*models.ColumnSchema{
							{Name: "id", DataType: "int", IsPrimaryKey: true},
							{Name: "name", DataType: "varchar(100)"},
						},
						PK: []string{"id"},
					},
				},
			},
		},
	}

	tgtSchema, err := c.ConvertSchema(srcSchema, "clickhouse")
	assert.NoError(t, err)
	assert.Len(t, tgtSchema.Databases, 1)
	assert.Len(t, tgtSchema.Databases[0].Tables, 1)

	tgtTable := tgtSchema.Databases[0].Tables[0]
	assert.Equal(t, "users", tgtTable.Name)
	assert.Contains(t, tgtTable.Engine, "MergeTree()")
	assert.Contains(t, tgtTable.Engine, "ORDER BY (id)")
	assert.Equal(t, "Int32", tgtTable.Columns[0].DataType)
}
