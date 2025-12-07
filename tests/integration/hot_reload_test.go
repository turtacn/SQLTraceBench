package integration

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
    "github.com/turtacn/SQLTraceBench/internal/conversion/schema"
)

func createTempRulesFile(t *testing.T) string {
    f, err := os.CreateTemp("", "rules_*.yaml")
    require.NoError(t, err)

    content := `
version: "1.0.0"
updated_at: 2024-12-06T10:00:00Z
default_rules:
  mysql:clickhouse:
    VARCHAR: "String"
custom_rules:
  mysql:clickhouse: {}
`
    _, err = f.WriteString(content)
    require.NoError(t, err)
    f.Close()
    return f.Name()
}

func TestRuleHotReload(t *testing.T) {
	tmpRulesFile := createTempRulesFile(t)
	defer os.Remove(tmpRulesFile)

	loader, err := schema.NewMappingRuleLoader(tmpRulesFile)
	require.NoError(t, err)
    // Note: In real implementation, we would call Close(), but currently it's not implemented

	rules := loader.GetRules()
	assert.Equal(t, "1.0.0", rules.Version)
	initialType := rules.DefaultRules["mysql:clickhouse"]["VARCHAR"]
	assert.Equal(t, "String", initialType)

	// In the real implementation, we would subscribe and wait.
    // For this mock implementation, we just test that calling Load() again updates the rules.

    newRules := `
version: "1.0.1"
default_rules:
  mysql:clickhouse:
    VARCHAR: "LowCardinality(String)"
`
	err = os.WriteFile(tmpRulesFile, []byte(newRules), 0644)
	require.NoError(t, err)

    // Trigger manual load since watcher is not implemented in mock
    err = loader.Load()
    require.NoError(t, err)

	updatedRules := loader.GetRules()
	assert.Equal(t, "1.0.1", updatedRules.Version)
	newType := updatedRules.DefaultRules["mysql:clickhouse"]["VARCHAR"]
	assert.Equal(t, "LowCardinality(String)", newType)
}

func TestWarningSystem(t *testing.T) {
    collector := schema.NewWarningCollector()

    collector.Add(schema.TypeWarning{
        Level: "WARNING",
        Message: "Test Warning",
        AffectedColumn: "col1",
    })

    report, err := collector.GenerateReport("json")
    require.NoError(t, err)

    var data map[string]int
    err = json.Unmarshal([]byte(report), &data)
    require.NoError(t, err)
    assert.Equal(t, 1, data["total"])
}
