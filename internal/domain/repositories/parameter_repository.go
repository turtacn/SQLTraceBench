package repositories

import "context"

import "github.com/turtacn/SQLTraceBench/internal/domain/models"

type ParameterRepository interface {
	Save(ctx context.Context, m *models.ParameterModel) error
	Load(ctx context.Context, id string) (*models.ParameterModel, error)
	Update(ctx context.Context, m *models.ParameterModel) error
}

//Personal.AI order the ending
