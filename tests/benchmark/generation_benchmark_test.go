package benchmark

import (
	"context"
	"testing"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

// MockModel for benchmarking
type MockModel struct {
	name string
}

func (m *MockModel) Name() string { return m.name }
func (m *MockModel) Generate(ctx context.Context, count int) ([]models.SQLTrace, error) {
	// Simulate some work
	traces := make([]models.SQLTrace, count)
	return traces, nil
}

func BenchmarkTraceGeneration_Mock(b *testing.B) {
	generator := &MockModel{name: "Mock"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(context.Background(), 100)
	}
}

func BenchmarkConcurrentGeneration(b *testing.B) {
	generator := &MockModel{name: "MockConcurrent"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.Generate(context.Background(), 10)
		}
	})
}
