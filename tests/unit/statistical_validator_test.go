package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestKSTest_IdenticalDistributions(t *testing.T) {
	// Generate identical distributions
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

	validator := services.NewStatisticalValidator(0.05, 0.01)
	result := validator.KolmogorovSmirnovTest(data, data)

	assert.True(t, result.Passed)
	assert.Equal(t, 0.0, result.Details["ks_statistic"]) // Statistic should be 0 for identical
	assert.Greater(t, result.PValue, 0.99)                // P-Value should be close to 1
}

func TestKSTest_DifferentDistributions(t *testing.T) {
	// Generate different distributions
	// Dist 1: Small numbers
	data1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	// Dist 2: Large numbers
	data2 := []float64{100.0, 200.0, 300.0, 400.0, 500.0}

	validator := services.NewStatisticalValidator(0.05, 0.01)
	result := validator.KolmogorovSmirnovTest(data1, data2)

	assert.False(t, result.Passed)
	assert.Less(t, result.PValue, 0.05)
}

func TestChiSquareTest_GoodnessOfFit(t *testing.T) {
	observed := []int{50, 30, 20}
	expected := []int{45, 35, 20} // Slight deviation

	validator := services.NewStatisticalValidator(0.05, 0.01)
	result := validator.ChiSquareTest(observed, expected)

	assert.True(t, result.Passed) // Deviation not significant
	assert.Greater(t, result.PValue, 0.05)
}

func TestChiSquareTest_SignificantDifference(t *testing.T) {
	observed := []int{10, 10, 80}
	expected := []int{33, 33, 33} // Uniform expected

	validator := services.NewStatisticalValidator(0.05, 0.01)
	result := validator.ChiSquareTest(observed, expected)

	assert.False(t, result.Passed)
	assert.Less(t, result.PValue, 0.01)
}
