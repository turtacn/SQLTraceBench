package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/turtacn/SQLTraceBench/internal/conversion/schema"
)

func TestPrecisionHandler_DecimalScale(t *testing.T) {
	h := schema.NewPrecisionHandler("")

	tests := []struct {
		name       string
		baseType   string
		params     []string
		targetDB   string
		expected   string
		expectLoss bool
	}{
		{
			name:     "Decimal32",
			baseType: "DECIMAL",
			params:   []string{"9", "2"},
			targetDB: "clickhouse",
			expected: "Decimal32(2)",
		},
		{
			name:     "Decimal64",
			baseType: "DECIMAL",
			params:   []string{"18", "4"},
			targetDB: "clickhouse",
			expected: "Decimal64(4)",
		},
		{
			name:     "Decimal128",
			baseType: "DECIMAL",
			params:   []string{"38", "10"},
			targetDB: "clickhouse",
			expected: "Decimal128(10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := h.HandlePrecision(tt.baseType, tt.params, tt.targetDB)
			assert.Equal(t, tt.expected, res.TargetType)
			assert.Equal(t, tt.expectLoss, res.HasLoss)
		})
	}
}

func TestPrecisionHandler_Timestamp(t *testing.T) {
	h := schema.NewPrecisionHandler("")

	tests := []struct {
		name       string
		baseType   string
		params     []string
		targetDB   string
		expected   string
	}{
		{
			name:     "No precision",
			baseType: "TIMESTAMP",
			params:   []string{},
			targetDB: "clickhouse",
			expected: "DateTime",
		},
        {
			name:     "Zero precision",
			baseType: "TIMESTAMP",
			params:   []string{"0"},
			targetDB: "clickhouse",
			expected: "DateTime",
		},
		{
			name:     "Millis",
			baseType: "TIMESTAMP",
			params:   []string{"3"},
			targetDB: "clickhouse",
			expected: "DateTime64(3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := h.HandlePrecision(tt.baseType, tt.params, tt.targetDB)
			assert.Equal(t, tt.expected, res.TargetType)
		})
	}
}
