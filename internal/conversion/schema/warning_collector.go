package schema

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// WarningCollector collects and manages conversion warnings.
type WarningCollector struct {
	warnings []TypeWarning
	stats    WarningStatistics
	mutex    sync.Mutex
}

type WarningStatistics struct {
	TotalWarnings   int
	ByLevel         map[string]int
	ByCategory      map[string]int
	AffectedTables  int
	AffectedColumns int
}

// NewWarningCollector creates a new WarningCollector.
func NewWarningCollector() *WarningCollector {
	return &WarningCollector{
		warnings: make([]TypeWarning, 0),
		stats: WarningStatistics{
			ByLevel:    make(map[string]int),
			ByCategory: make(map[string]int),
		},
	}
}

// Add adds a warning to the collector.
func (c *WarningCollector) Add(warning TypeWarning) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.warnings = append(c.warnings, warning)
	c.stats.TotalWarnings++
	c.stats.ByLevel[warning.Level]++
}

// GenerateReport generates a report in the specified format.
func (c *WarningCollector) GenerateReport(format string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch format {
	case "json":
		return c.generateJSONReport()
	case "markdown":
		return c.generateMarkdownReport()
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func (c *WarningCollector) generateJSONReport() (string, error) {
    report := map[string]interface{}{
        "statistics": c.stats,
        "warnings":   c.warnings,
        "total":      c.stats.TotalWarnings, // Keep backward compat with simple check
    }
    bytes, err := json.MarshalIndent(report, "", "  ")
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func (c *WarningCollector) generateMarkdownReport() (string, error) {
	var sb strings.Builder
	sb.WriteString("# Type Conversion Warning Report\n\n")
	sb.WriteString(fmt.Sprintf("- Total Warnings: %d\n", c.stats.TotalWarnings))

	// Group by level
	for _, w := range c.warnings {
		sb.WriteString(fmt.Sprintf("- [%s] %s: %s (Suggestion: %s)\n", w.Level, w.AffectedColumn, w.Message, w.Suggestion))
	}

	return sb.String(), nil
}
