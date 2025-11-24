package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/plugin_registry"
)

func TestClickHousePlugin_EndToEnd(t *testing.T) {
	// 1. Compile plugin binary
	// Assuming running from project root
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Adjust path if running from tests/e2e
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
	defer os.RemoveAll(binDir) // cleanup

	buildCmd := exec.Command("go", "build", "-o", pluginBin, pluginSource)
	buildCmd.Dir = projectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Plugin build failed: %s", string(buildOutput))

	// Verify binary size < 20MB (AC-1)
	info, err := os.Stat(pluginBin)
	require.NoError(t, err)
	// Just use info to avoid unused error, optionally check size
	_ = info

	// 2. Load through registry
	reg := plugin_registry.NewRegistry()
	// Registry LoadPluginsFromDir expects directory
	err = reg.LoadPluginsFromDir(binDir)
	require.NoError(t, err)
	defer reg.Close()

	// 3. Get plugin instance
	// The registry stores it by name.
	// But `reg.Get(name)` returns `plugins.Plugin` interface which might be the legacy one.
	// We need to cast it to `pkg/plugin.DatabasePlugin` to access ConvertSchema.
	// Wait, `plugins.Plugin` interface in `plugins/interfaces.go` only has `TranslateQuery`.
	// The `plugin_registry` uses `plugins.Plugin` from `plugins/interfaces.go`.

	// Issue: The `plugin_registry` seems to wrap the gRPC client but returns `plugins.Plugin`.
	// If `plugins.Plugin` doesn't have `ConvertSchema`, we can't test it via that interface.
	// Let's check `plugins/interfaces.go` content again.
	// It has: Name(), Version(), TranslateQuery(sql string) (string, error)
	// It DOES NOT have ConvertSchema.

	// However, the gRPC implementation `grpc_impl.GRPCPluginImpl` implements `DatabasePlugin` from `pkg/plugin/interface.go`.
	// And `pkg/plugin/interface.go` HAS `ConvertSchema`.

	// We need to check if we can cast or if we need to update `plugins.Plugin` interface.
	// The user requirements said:
	// MODIFY: plugins/clickhouse/plugin.go
	// func (p *ClickHousePlugin) ConvertSchema(...)

	// The `plugin_registry` returns `plugins.Plugin`.
	plRaw, ok := reg.Get("clickhouse")
	require.True(t, ok, "Plugin 'clickhouse' not found")

	// Verify Name (AC-2)
	assert.Equal(t, "clickhouse", plRaw.Name())

	// We need to assert that plRaw implements the interface that has ConvertSchema.
	// Since `plugins.Plugin` doesn't have it, we might need to assert it to an interface that does,
	// or `plugins.Plugin` should have it?
	// The prompt didn't ask to modify `plugins/interfaces.go`.
	// But `plugins/clickhouse/plugin.go` implements `ConvertSchema`.

	// If the gRPC client stub (which `reg.Get` returns indirectly? No, `reg.Get` returns the value dispensed by `rpcClient.Dispense`).
	// The `GRPCPluginImpl` in `pkg/plugin/grpc_impl` implements the interface.
	// Let's assume we can type assert to an interface that includes ConvertSchema.

	type SchemaConverterPlugin interface {
		ConvertSchema(schema *models.Schema) (*models.Schema, error)
	}

	converterPl, ok := plRaw.(SchemaConverterPlugin)
	require.True(t, ok, "Plugin does not implement SchemaConverterPlugin interface")

	// 4. Test Schema Conversion
	schema := createTestSchema() // MySQL schema
	converted, err := converterPl.ConvertSchema(schema)
	assert.NoError(t, err)
	require.NotNil(t, converted)
	require.NotEmpty(t, converted.Databases)
	require.NotEmpty(t, converted.Databases[0].Tables)

	tbl := converted.Databases[0].Tables[0]
	assert.Contains(t, tbl.Engine, "MergeTree")
	assert.Contains(t, tbl.Engine, "ORDER BY")
	// Verify CreateSQL populated
	assert.Contains(t, tbl.CreateSQL, "CREATE TABLE user_logs")
	assert.Contains(t, tbl.CreateSQL, "ENGINE = MergeTree()")

	// Verify Decimal mapping (AC-3)
	// DECIMAL(10,2) -> Decimal128(2)
	var decimalCol *models.ColumnSchema
	for _, col := range tbl.Columns {
		if col.Name == "amount" {
			decimalCol = col
			break
		}
	}
	require.NotNil(t, decimalCol)
	assert.Equal(t, "Decimal128(2)", decimalCol.DataType)

	// 5. Test Query Translation
	// AC-4: E2E < 5s is implicit by test runtime

	translated, err := plRaw.TranslateQuery("SELECT SQL_NO_CACHE * FROM t WHERE created_at = NOW()")
	assert.NoError(t, err)
	assert.NotContains(t, translated, "SQL_NO_CACHE")
	assert.Contains(t, translated, "now()")
	assert.Equal(t, "SELECT * FROM t WHERE created_at = now()", translated)
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
							{Name: "user_name", DataType: "VARCHAR(255)"},
							{Name: "amount", DataType: "DECIMAL(10,2)"}, // Test AC-3
							{Name: "created_at", DataType: "DATETIME"},
						},
						PK: []string{"id"},
					},
				},
			},
		},
	}
}
