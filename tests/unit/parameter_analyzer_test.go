package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestParameterAnalyzer_TypeInference(t *testing.T) {
	traces := []models.SQLTrace{
		{Parameters: map[string]interface{}{"id": "123", "name": "Alice", "date": "2025-01-01"}},
		{Parameters: map[string]interface{}{"id": 456, "name": "Bob", "date": "2025-01-02T15:04:05Z"}},
	}

	analyzer := services.NewParameterAnalyzer()
	stats := analyzer.Analyze(traces)

	assert.Equal(t, "INT", stats["id"].DataType)
	assert.Equal(t, "STRING", stats["name"].DataType)
	assert.Equal(t, "DATETIME", stats["date"].DataType)
	assert.Equal(t, 2, stats["id"].Cardinality)
}
