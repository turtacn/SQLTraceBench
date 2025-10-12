package services

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaLoader defines the interface for loading a database schema.
type SchemaLoader interface {
	// LoadSchema connects to a database and applies the given schema.
	LoadSchema(ctx context.Context, schema *models.DatabaseSchema) error
}