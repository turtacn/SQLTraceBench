package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateQuery(t *testing.T) {
	translator := NewQueryTranslator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Keep Backticks",
			input:    "SELECT * FROM `User`",
			expected: "SELECT * FROM `User`",
		},
		{
			name:     "Remove SQL_NO_CACHE",
			input:    "SELECT SQL_NO_CACHE * FROM t",
			expected: "SELECT * FROM t",
		},
		{
			name:     "Normalize Function Name",
			input:    "SELECT NOW()",
			expected: "SELECT now()",
		},
		{
			name:     "Remove Semicolon",
			input:    "SELECT * FROM t;",
			expected: "SELECT * FROM t",
		},
		{
			name:     "Combined",
			input:    "SELECT SQL_NO_CACHE * FROM `t` WHERE x = NOW();",
			expected: "SELECT * FROM `t` WHERE x = now()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translator.TranslateQuery(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
