package starrocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestTypeMapping(t *testing.T) {
	c := &StarRocksConverter{}

	tests := []struct {
		input    string
		expected string
	}{
		{"INT(11)", "INT"},
		{"VARCHAR(255)", "VARCHAR(255)"},
		{"TEXT", "STRING"},
		{"TINYINT(1)", "TINYINT(1)"},
		{"MEDIUMINT", "INT"},
		{"TIMESTAMP", "DATETIME"},
		{"FLOAT", "FLOAT"},
		{"DOUBLE", "DOUBLE"},
		{"BOOLEAN", "BOOLEAN"},
	}

	for _, test := range tests {
		result := c.mapType(test.input)
		assert.Equal(t, test.expected, result, "Type mapping failed for %s", test.input)
	}
}

func TestDistributedKey(t *testing.T) {
	c := &StarRocksConverter{}

	// Case 1: With PK
	sourceTblWithPK := &models.TableSchema{
		Name:    "users",
		PK:      []string{"id"},
		Columns: []*models.ColumnSchema{{Name: "id", DataType: "INT"}, {Name: "name", DataType: "VARCHAR(100)"}},
	}
	resWithPK := c.convertTable(sourceTblWithPK)
	assert.Contains(t, resWithPK.Engine, "OLAP", "Engine should be OLAP")
	// Now with backticks
	assert.Contains(t, resWithPK.Engine, "DUPLICATE KEY(`id`)", "Should contain DUPLICATE KEY")
	assert.Contains(t, resWithPK.Engine, "DISTRIBUTED BY HASH(`id`)", "Should contain DISTRIBUTED BY HASH(id)")

	// Case 2: No PK
	sourceTblNoPK := &models.TableSchema{
		Name:    "logs",
		Columns: []*models.ColumnSchema{{Name: "log_id", DataType: "INT"}, {Name: "msg", DataType: "TEXT"}},
	}
	resNoPK := c.convertTable(sourceTblNoPK)
	assert.Contains(t, resNoPK.Engine, "OLAP", "Engine should be OLAP")
	// If no PK, we used fallback to first column for distribution in our implementation
	// Now with backticks
	assert.Contains(t, resNoPK.Engine, "DUPLICATE KEY(`log_id`)", "Should contain DUPLICATE KEY using fallback")
	assert.Contains(t, resNoPK.Engine, "DISTRIBUTED BY HASH(`log_id`)", "Should contain DISTRIBUTED BY HASH using fallback")
}

func TestSchemaConversion(t *testing.T) {
	c := &StarRocksConverter{}

	sourceSchema := &models.Schema{
		Databases: []models.DatabaseSchema{
			{
				Name: "testdb",
				Tables: []*models.TableSchema{
					{
						Name: "t1",
						PK: []string{"id"},
						Columns: []*models.ColumnSchema{
							{Name: "id", DataType: "INT(11)"},
						},
					},
				},
			},
		},
	}

	res, err := c.ConvertSchema(sourceSchema)
	assert.NoError(t, err)
	assert.Equal(t, "testdb", res.Databases[0].Name)
	assert.Equal(t, "t1", res.Databases[0].Tables[0].Name)
	assert.Equal(t, "INT", res.Databases[0].Tables[0].Columns[0].DataType)
}
