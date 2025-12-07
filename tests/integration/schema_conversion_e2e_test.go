package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

func TestSchemaConversion_E2E(t *testing.T) {
	// Create temporary directory for inputs and outputs
	tmpDir, err := os.MkdirTemp("", "schema_conversion_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Define test cases
	testCases := []struct {
		name       string
		sourceDB   string
		targetDB   string
		inputDDL   string
		expectFile bool
	}{
		{
			name:     "MySQL to ClickHouse",
			sourceDB: "mysql",
			targetDB: "clickhouse",
			inputDDL: `
CREATE TABLE customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT NOW()
);
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT,
    amount DECIMAL(10,2),
    status ENUM('pending', 'shipped')
);
`,
			expectFile: true,
		},
		{
			name:     "PostgreSQL to ClickHouse",
			sourceDB: "postgres",
			targetDB: "clickhouse",
			inputDDL: `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    data JSONB,
    tags INTEGER[]
);
`,
			expectFile: true,
		},
		{
			name:     "TiDB to ClickHouse",
			sourceDB: "tidb",
			targetDB: "clickhouse",
			inputDDL: `
CREATE TABLE logs (
    id BIGINT AUTO_RANDOM PRIMARY KEY,
    msg TEXT
) SHARD_ROW_ID_BITS=4;
`,
			expectFile: true,
		},
	}

	parser := parsers.NewStreamingTraceParser(0)
	svc := conversion.NewService(parser, plugin_registry.GlobalRegistry)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputPath := filepath.Join(tmpDir, tc.name+"_input.sql")
			outputPath := filepath.Join(tmpDir, tc.name+"_output.sql")

			err := os.WriteFile(inputPath, []byte(tc.inputDDL), 0644)
			assert.NoError(t, err)

			// Mock request
			req := conversion.ConvertRequest{
				SourceSchemaPath: inputPath,
				TargetDBType:     tc.targetDB,
				OutputPath:       outputPath,
				SourceDB:         tc.sourceDB,
			}

			err = svc.ConvertSchemaFromFile(context.Background(), req)
			assert.NoError(t, err)

			if tc.expectFile {
				assert.FileExists(t, outputPath)
				content, err := os.ReadFile(outputPath)
				assert.NoError(t, err)
				t.Logf("Converted SQL for %s:\n%s", tc.name, string(content))

				// Basic validations
				sContent := string(content)
				if tc.targetDB == "clickhouse" {
					assert.Contains(t, sContent, "MergeTree")
				}
			}
		})
	}
}
