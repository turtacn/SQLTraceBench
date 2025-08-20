package repositories

import (
	"context"
	"time"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type TraceRepository interface {
	Save(ctx context.Context, trace models.SQLTrace) error
	SaveBatch(ctx context.Context, traces []models.SQLTrace) error
	LoadByRange(ctx context.Context, from, to time.Time) ([]models.SQLTrace, error)
	ListTables() ([]string, error)
	Count(ctx context.Context) (int64, error)
}

type traceRepo struct{}

func NewTraceRepository() TraceRepository {
	return &traceRepo{}
}

func (r *traceRepo) Save(ctx context.Context, trace models.SQLTrace) error {
	return nil // TODO implement
}

func (r *traceRepo) SaveBatch(ctx context.Context, traces []models.SQLTrace) error {
	return nil
}

func (r *traceRepo) LoadByRange(ctx context.Context, from, to time.Time) ([]models.SQLTrace, error) {
	return nil, nil
}

func (r *traceRepo) ListTables() ([]string, error) {
	return nil, nil
}

func (r *traceRepo) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

//Personal.AI order the ending
