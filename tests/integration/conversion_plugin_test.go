package integration

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"github.com/turtacn/SQLTraceBench/plugins/starrocks"
)

func TestConversionWithStarRocksPlugin(t *testing.T) {
	// 1. Setup
	tempDir, err := ioutil.TempDir("", "conversion_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	sourcePath := filepath.Join(tempDir, "source.sql")
	outputPath := filepath.Join(tempDir, "output.sql")

	// Create a dummy source schema file (MySQL dialect)
	sourceSQL := `
CREATE TABLE users (
	id INT PRIMARY KEY,
	name VARCHAR(255),
	created_at DATETIME
);
`
	err = ioutil.WriteFile(sourcePath, []byte(sourceSQL), 0644)
	require.NoError(t, err)

	// 2. Register Plugin Manually
	srPlugin := starrocks.New()
	reg := plugin_registry.NewRegistry()
	reg.Register(srPlugin)

	// 3. Create Service
	parser := parsers.NewRegexParser() // Simple parser
	svc := conversion.NewService(parser, reg)

	// 4. Execute Conversion
	req := conversion.ConvertRequest{
		SourceSchemaPath: sourcePath,
		TargetDBType:     "starrocks",
		OutputPath:       outputPath,
	}
	err = svc.ConvertSchemaFromFile(context.Background(), req)
	require.NoError(t, err)

	// 5. Verify Output
	content, err := ioutil.ReadFile(outputPath)
	require.NoError(t, err)
	outputSQL := string(content)

	// Verification logic based on StarRocks plugin implementation
	// We expect: ENGINE=OLAP ... DISTRIBUTED BY HASH(`id`) ...
	assert.Contains(t, outputSQL, "CREATE TABLE `users`")
	assert.Contains(t, outputSQL, "ENGINE=OLAP")
	assert.Contains(t, outputSQL, "DISTRIBUTED BY HASH(`id`)")
	assert.Contains(t, outputSQL, "DUPLICATE KEY(`id`)")
	assert.Contains(t, outputSQL, "`name` VARCHAR(255)")
}
