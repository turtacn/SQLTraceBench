package repositories

import "context"

import "github.com/turtacn/SQLTraceBench/internal/domain/models"

type ParameterRepository interface {
	Save(ctx context.Context, m *models.WorkloadParameterModel) error
	Load(ctx context.Context, id string) (*models.WorkloadParameterModel, error)
	Update(ctx context.Context, m *models.WorkloadParameterModel) error
}
