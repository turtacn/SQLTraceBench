package repositories

import (
	"context"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type TemplateRepository interface {
	Save(ctx context.Context, tpl models.SQLTemplate) error
	List(ctx context.Context) ([]models.SQLTemplate, error)
	FindByTable(ctx context.Context, table string) ([]models.SQLTemplate, error)
	Delete(ctx context.Context, id string) error
}

//Personal.AI order the ending
