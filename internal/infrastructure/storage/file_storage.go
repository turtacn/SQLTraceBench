package storage

import (
	"errors"
	"os"

	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

type FileStorage struct {
	baseDir string
}

func NewFileStorage(base string) (*FileStorage, error) {
	if base == "" {
		base = os.TempDir()
	}
	return &FileStorage{baseDir: base}, nil
}

func (fs *FileStorage) SaveWorkload(w *models.BenchmarkWorkload) error {
	return errors.New("not implemented for MVP")
}

//Personal.AI order the ending
