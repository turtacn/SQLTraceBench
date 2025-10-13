package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerformanceMetrics_QPS(t *testing.T) {
	pm := &PerformanceMetrics{
		QueriesExecuted: 100,
		Duration:        10 * time.Second,
	}
	assert.Equal(t, 10.0, pm.QPS())
}

func TestPerformanceMetrics_ErrorRate(t *testing.T) {
	pm := &PerformanceMetrics{
		QueriesExecuted: 100,
		Errors:          10,
	}
	assert.Equal(t, 0.1, pm.ErrorRate())
}

func TestPerformanceMetrics_CalculatePercentiles(t *testing.T) {
	pm := &PerformanceMetrics{
		Latencies: []time.Duration{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	pm.CalculatePercentiles()
	assert.Equal(t, time.Duration(6), pm.P50)
	assert.Equal(t, time.Duration(10), pm.P90)
	assert.Equal(t, time.Duration(10), pm.P99)
}

func TestSQLTemplate_ExtractParameters(t *testing.T) {
	tpl := &SQLTemplate{
		RawSQL: "SELECT * FROM users WHERE id = :id AND name = :name",
	}
	tpl.ExtractParameters()
	assert.Equal(t, []string{":id", ":name"}, tpl.Parameters)
}

func TestSQLTemplate_GenerateQuery(t *testing.T) {
	tpl := &SQLTemplate{
		RawSQL:     "SELECT * FROM users WHERE id = :id AND name = :name",
		Parameters: []string{":id", ":name"},
	}
	params := map[string]interface{}{
		":id":   1,
		":name": "test",
	}
	q, err := tpl.GenerateQuery(params)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM users WHERE id = ? AND name = ?", q.Query)
	assert.Equal(t, []interface{}{1, "test"}, q.Args)
}