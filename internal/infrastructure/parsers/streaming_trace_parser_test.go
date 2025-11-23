package parsers

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestStreamingParser_ValidJSON(t *testing.T) {
	jsonl := `{"timestamp":"2025-01-01T00:00:00Z","query_text":"SELECT 1","latency":0.1}
{"timestamp":"2025-01-01T00:00:01Z","query_text":"SELECT 2","latency":0.2}`

	reader := strings.NewReader(jsonl)
	parser := StreamingTraceParser{}

	var count int
	err := parser.Parse(reader, func(trace models.SQLTrace) error {
		count++
		if count == 1 {
			assert.Equal(t, "SELECT 1", trace.Query)
			assert.Equal(t, 100*time.Millisecond, trace.Latency)
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestStreamingParser_CorruptedLine(t *testing.T) {
	jsonl := `{"timestamp":"2025-01-01T00:00:00Z","query_text":"SELECT 1"}
INVALID_JSON
{"timestamp":"2025-01-01T00:00:02Z","query_text":"SELECT 2"}`

	reader := strings.NewReader(jsonl)
	parser := StreamingTraceParser{}

	var count int
	err := parser.Parse(reader, func(trace models.SQLTrace) error {
		count++
		return nil
	})

	// Should not error out, just skip the invalid line
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestStreamingParser_InvalidTimestamp(t *testing.T) {
	jsonl := `{"timestamp":"INVALID_DATE","query_text":"SELECT 1"}`

	reader := strings.NewReader(jsonl)
	parser := StreamingTraceParser{}

	var count int
	err := parser.Parse(reader, func(trace models.SQLTrace) error {
		count++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func BenchmarkStreamingParser_MemoryStability(b *testing.B) {
	// Generate a repeatable large input
	line := `{"timestamp":"2025-01-01T00:00:00Z","query_text":"SELECT * FROM table WHERE id = 1","latency":0.001}` + "\n"

	// Create a reader that repeats this line
	// For benchmark, we might want to process N lines.
	// But b.N is loop count.

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(line)
		parser := StreamingTraceParser{}
		_ = parser.Parse(reader, func(trace models.SQLTrace) error {
			return nil
		})
	}
}
