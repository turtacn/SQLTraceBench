package docs_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper to extract code blocks from markdown
func extractCodeBlocks(filepath string, language string) ([]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Regex for finding code blocks: ```language ... ```
	re := regexp.MustCompile("(?s)```" + language + "\\s+(.*?)```")
	matches := re.FindAllStringSubmatch(string(content), -1)

	var blocks []string
	for _, match := range matches {
		if len(match) > 1 {
			blocks = append(blocks, strings.TrimSpace(match[1]))
		}
	}
	return blocks, nil
}

// TestDocLinks checks if internal links in markdown files are valid
// Ideally this would check file existence, but for now we just parse them
func TestQuickstartExample(t *testing.T) {
	// Find the file relative to the repo root.
	// Note: tests usually run in the directory they are in, so we might need to go up.
	path := "../../docs/user_guide/quickstart.md"

	// Check if file exists
	_, err := os.Stat(path)
	if err != nil {
		t.Skipf("Quickstart doc not found at %s", path)
	}

	blocks, err := extractCodeBlocks(path, "bash")
	assert.NoError(t, err)

	// Basic assertion that we found some bash blocks
	assert.Greater(t, len(blocks), 0, "Should find bash code blocks in quickstart guide")

	// We can't easily 'run' them without a full environment, but we can syntax check or look for expected commands
	foundGenerate := false
	for _, block := range blocks {
		if strings.Contains(block, "sql_trace_bench generate") {
			foundGenerate = true
			break
		}
	}
	assert.True(t, foundGenerate, "Quickstart should contain a generate command")
}

func TestArchitectureDocExistence(t *testing.T) {
	files := []string{
		"../../docs/architecture/system_architecture.md",
		"../../docs/architecture/data_flow.md",
	}

	for _, f := range files {
		_, err := os.Stat(f)
		assert.NoError(t, err, "File %s should exist", f)
	}
}

func TestAPIDocExistence(t *testing.T) {
    files := []string{
        "../../docs/api/rest_api.md",
        "../../docs/api/cli_reference.md",
    }
    for _, f := range files {
        _, err := os.Stat(f)
        assert.NoError(t, err, "File %s should exist", f)
    }
}
