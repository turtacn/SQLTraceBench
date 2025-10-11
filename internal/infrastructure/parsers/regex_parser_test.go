package parsers

import (
	"reflect"
	"testing"
)

func TestRegexParser_ListTables(t *testing.T) {
	parser := NewRegexParser()

	testCases := []struct {
		name     string
		sql      string
		expected []string
	}{
		{
			name:     "simple select",
			sql:      "select * from users",
			expected: []string{"users"},
		},
		{
			name:     "multiple joins",
			sql:      "select u.id, o.order_id from users u join orders o on u.id = o.user_id",
			expected: []string{"users", "orders"},
		},
		{
			name:     "mixed case keywords",
			sql:      "SELECT * FrOm customers JOIN sales ON customers.id = sales.customer_id",
			expected: []string{"customers", "sales"},
		},
		{
			name:     "schema qualified tables",
			sql:      "select * from public.users join internal.events on users.id = events.user_id",
			expected: []string{"public.users", "internal.events"},
		},
		{
			name:     "no tables",
			sql:      "select 1 + 1",
			expected: []string{},
		},
		{
			name:     "deduplicate tables",
			sql:      "select * from users join users_metadata on users.id = users_metadata.user_id",
			expected: []string{"users", "users_metadata"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tables, err := parser.ListTables(tc.sql)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tables, tc.expected) {
				t.Errorf("expected tables %v, but got %v", tc.expected, tables)
			}
		})
	}
}