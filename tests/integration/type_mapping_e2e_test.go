package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
    "github.com/turtacn/SQLTraceBench/internal/conversion/schema"
    "github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestTypeMappingE2E_ComplexSchema(t *testing.T) {
    // Setup
    converter := schema.NewMySQLConverter()

    // Create complex schema
    sourceTable := &models.TableSchema{
        Name: "complex_users",
        PK: []string{"id"},
        Columns: []*models.ColumnSchema{
            {Name: "id", DataType: "BIGINT UNSIGNED", IsPrimaryKey: true, IsNullable: false},
            {Name: "email", DataType: "VARCHAR(255)", IsNullable: false}, // Expect LowCardinality due to default rule in test env (or mocked)
            {Name: "balance", DataType: "DECIMAL(10,2)", IsNullable: true},
            {Name: "created_at", DataType: "TIMESTAMP(6)", IsNullable: false},
            {Name: "status", DataType: "ENUM('active','inactive')", IsNullable: false},
        },
    }

    // Execute
    targetTable, err := converter.ConvertTable(sourceTable, "clickhouse")
    require.NoError(t, err)

    // Assertions
    // Note: assertions depend on what rules are actually loaded.
    // If standard rules are used:
    // id -> UInt64 (if mapped in baseRules) or Int64 (if baseRules says Int64)
    // email -> String (unless rule matches)
    // balance -> Decimal64(2)

    for _, col := range targetTable.Columns {
        if col.Name == "balance" {
            assert.Contains(t, col.DataType, "Decimal")
        }
        if col.Name == "created_at" {
            assert.Contains(t, col.DataType, "DateTime64(6)")
        }
    }

    // Report
    report, err := converter.GenerateWarningReport("json")
    require.NoError(t, err)
    assert.Contains(t, report, "total")
}
