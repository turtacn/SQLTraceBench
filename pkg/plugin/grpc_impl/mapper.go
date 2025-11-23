package grpc_impl

import (
	"encoding/json"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// ToProtoSchema converts a domain Schema to its proto string representation (JSON).
func ToProtoSchema(s *models.Schema) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromProtoSchema converts a proto string representation (JSON) back to a domain Schema.
func FromProtoSchema(data string) (*models.Schema, error) {
	var s models.Schema
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// TraceToQuery extracts the query string from a domain SQLTrace.
// This is a helper as SQLTrace <-> proto.SQLTrace (which is query string in TranslateQuery)
func TraceToQuery(t *models.SQLTrace) string {
	if t == nil {
		return ""
	}
	return t.Query
}

// QueryToTrace creates a minimal SQLTrace from a query string.
func QueryToTrace(query string) *models.SQLTrace {
	return &models.SQLTrace{
		Query: query,
	}
}
