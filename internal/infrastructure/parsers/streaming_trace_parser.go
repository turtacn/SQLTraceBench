package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
)

// StreamingTraceParser is responsible for parsing trace files in a streaming fashion.
type StreamingTraceParser struct {
	BufferSize int
}

// NewStreamingTraceParser creates a new StreamingTraceParser with the given buffer size.
func NewStreamingTraceParser(bufferSize int) *StreamingTraceParser {
	if bufferSize <= 0 {
		bufferSize = 1024 * 1024 // Default 1MB
	}
	return &StreamingTraceParser{
		BufferSize: bufferSize,
	}
}

// traceDTO is a data transfer object used for unmarshalling the JSON lines.
// It handles the mapping between the JSON fields and the models.SQLTrace struct.
type traceDTO struct {
	Query     string  `json:"query_text"`
	QueryAlt  string  `json:"query"` // Fallback
	Timestamp string  `json:"timestamp"`
	Latency   float64 `json:"latency"`
}

// Parse reads the provided reader line by line, unmarshals each line into a SQLTrace object,
// and calls the callback function with the parsed trace.
func (p *StreamingTraceParser) Parse(reader io.Reader, callback func(models.SQLTrace) error) error {
	scanner := bufio.NewScanner(reader)
	// Set buffer limit based on configuration
	buf := make([]byte, 64*1024)
	bufferSize := p.BufferSize
	if bufferSize == 0 {
		bufferSize = 1024 * 1024 // Default fallback
	}
	scanner.Buffer(buf, bufferSize)

	lineNum := 0
	logger := utils.GetGlobalLogger()

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}

		var dto traceDTO
		if err := json.Unmarshal(line, &dto); err != nil {
			logger.Error("Skipped malformed line", utils.Field{Key: "line", Value: lineNum}, utils.Field{Key: "error", Value: err})
			continue
		}

		// Map DTO to domain model
		// Handle timestamp parsing (assuming ISO8601/RFC3339)
		ts, err := time.Parse(time.RFC3339, dto.Timestamp)
		if err != nil {
			// Try fallback format if RFC3339 fails, e.g. "2025-01-01 00:00:00"
			// For now, log error and skip
			logger.Error("Invalid timestamp format", utils.Field{Key: "line", Value: lineNum}, utils.Field{Key: "timestamp", Value: dto.Timestamp})
			continue
		}

		query := dto.Query
		if query == "" {
			query = dto.QueryAlt
		}

		trace := models.SQLTrace{
			Query:     query,
			Timestamp: ts,
			// Assuming Latency is in seconds if float. If it's ms, use Millisecond.
			// I'll use Second as default for float latency.
			Latency:   time.Duration(dto.Latency * float64(time.Second)),
		}

		if err := callback(trace); err != nil {
			return fmt.Errorf("callback failed at line %d: %w", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}
