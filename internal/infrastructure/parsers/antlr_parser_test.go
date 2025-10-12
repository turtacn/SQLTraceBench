package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAntlrParser_ListTables(t *testing.T) {
	parser := NewAntlrParser()

	testCases := []struct {
		name     string
		sql      string
		expected []string
	}{
		{
			name:     "simple select",
			sql:      "SELECT * FROM users;",
			expected: []string{"users"},
		},
		{
			name:     "select with where clause",
			sql:      "SELECT * FROM customers WHERE id = 1;",
			expected: []string{"customers"},
		},
		{
			name:     "multiple joins",
			sql:      "SELECT u.id, o.order_id FROM users JOIN orders ON u.id = o.user_id;",
			expected: []string{"users", "orders"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tables, err := parser.ListTables(tc.sql)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, tables, "extracted tables should match expected tables")
		})
	}
}