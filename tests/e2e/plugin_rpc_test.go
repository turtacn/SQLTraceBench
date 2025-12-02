package e2e

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
	"github.com/turtacn/SQLTraceBench/pkg/proto"
)

func TestClickHousePlugin_EndToEnd(t *testing.T) {
	// 1. Compile plugin binary
	cwd, err := os.Getwd()
	require.NoError(t, err)

	projectRoot := cwd
	if filepath.Base(cwd) == "e2e" {
		projectRoot = filepath.Join(cwd, "../..")
	} else if filepath.Base(cwd) == "tests" {
		projectRoot = filepath.Join(cwd, "..")
	}

	pluginSource := filepath.Join(projectRoot, "cmd/plugins/clickhouse/main.go")
	binDir := filepath.Join(projectRoot, "bin")
	pluginBin := filepath.Join(binDir, "clickhouse-plugin")

	err = os.MkdirAll(binDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(binDir)

	buildCmd := exec.Command("go", "build", "-o", pluginBin, pluginSource)
	buildCmd.Dir = projectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Plugin build failed: %s", string(buildOutput))

	// 2. Load through registry
	reg := plugin_registry.NewRegistry()
	err = reg.LoadPluginsFromDir(binDir)
	require.NoError(t, err)
	defer reg.Close()

	// 3. Get plugin instance
	plRaw, ok := reg.Get("clickhouse")
	require.True(t, ok, "Plugin 'clickhouse' not found")

	// 4. Test Schema Conversion
	schema := createTestSchema()
	converted, err := plRaw.ConvertSchema(schema)
	assert.NoError(t, err)
	require.NotNil(t, converted)

	// 5. Test Query Translation
	translated, err := plRaw.TranslateQuery("SELECT SQL_NO_CACHE * FROM t WHERE created_at = NOW()")
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM t WHERE created_at = now()", translated)

	// 6. Test Query Execution
	_, err = plRaw.ExecuteQuery(context.Background(), &proto.ExecuteQueryRequest{Sql: "SELECT 1"})
	assert.NoError(t, err)
}

func createTestSchema() *models.Schema {
	return &models.Schema{
		Databases: []models.DatabaseSchema{
			{
				Name: "test_db",
				Tables: []*models.TableSchema{
					{
						Name: "user_logs",
						Columns: []*models.ColumnSchema{
							{Name: "id", DataType: "BIGINT", IsPrimaryKey: true},
						},
						PK: []string{"id"},
					},
				},
			},
		},
	}
}
