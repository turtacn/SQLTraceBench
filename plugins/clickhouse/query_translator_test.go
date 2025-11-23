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
			name:     "Remove Backticks",
			input:    "SELECT * FROM `User`",
			expected: "SELECT * FROM \"User\"",
		},
		{
			name:     "Remove Engine",
			input:    "CREATE TABLE t (id int) ENGINE=InnoDB",
			expected: "CREATE TABLE t (id int)",
		},
		{
			name:     "Remove Charset",
			input:    "CREATE TABLE t (id int) DEFAULT CHARSET=utf8mb4",
			expected: "CREATE TABLE t (id int)",
		},
		{
			name:     "Remove AutoIncrement",
			input:    "id INT AUTO_INCREMENT",
			expected: "id INT",
		},
		{
			name:     "Complex",
			input:    "CREATE TABLE `Orders` (`id` int) ENGINE=InnoDB DEFAULT CHARSET=utf8;",
			expected: "CREATE TABLE \"Orders\" (\"id\" int);",
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
