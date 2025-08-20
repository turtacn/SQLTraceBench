package services

import (
	"context"
	"strings"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type SchemaService interface {
	Convert(ctx context.Context, srcSchema *models.DatabaseSchema, toDialect string) (*models.DatabaseSchema, error)
}

type DefaultSchemaService struct{}

func NewSchemaService() SchemaService { return &DefaultSchemaService{} }

func (s *DefaultSchemaService) Convert(_ context.Context, src *models.DatabaseSchema, toDialect string) (*models.DatabaseSchema, error) {
	// minimal: rename keywords -> lower case + engine string swap
	dst := &models.DatabaseSchema{
		DatabaseName: src.DatabaseName,
	}
	for _, tbl := range src.TableDefinitions {
		newTbl := tbl
		newTbl.Engine = strings.ReplaceAll(tbl.Engine, "MyISAM", "MergeTree") // demo
		dst.TableDefinitions = append(dst.TableDefinitions, newTbl)
	}
	return dst, nil
}
