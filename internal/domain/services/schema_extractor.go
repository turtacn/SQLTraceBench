package services

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// SchemaExtractor defines the interface for extracting a database schema.
type SchemaExtractor interface {
	// ExtractSchema connects to a database and extracts its schema.
	ExtractSchema(ctx context.Context) (*models.DatabaseSchema, error)
}