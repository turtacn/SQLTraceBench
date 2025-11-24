package services_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestSynthesizer_FillParameters(t *testing.T) {
	// Setup Parameter Models
	// Template 1: SELECT * FROM users WHERE id = :id AND region = :region

	// id: Zipfian
	idModel := &models.ParameterModel{
		ParamName:   ":id",
		DataType:    "INT",
		DistType:    models.DistZipfian,
		ZipfS:       2.0, // Highly skewed
		Cardinality: 3,
		TopValues:   []interface{}{100, 200, 300}, // Rank 0 -> 100, Rank 1 -> 200, ...
	}

	// region: Uniform/Weighted
	regionModel := &models.ParameterModel{
		ParamName:      ":region",
		DataType:       "STRING",
		DistType:       models.DistUniform,
		Cardinality:    2,
		TopValues:      []interface{}{"US", "EU"},
		TopFrequencies: []int{50, 50},
	}

	wlParams := &models.WorkloadParameterModel{
		TemplateParameters: map[string]map[string]*models.ParameterModel{
			"group1": {
				":id":     idModel,
				":region": regionModel,
			},
		},
	}

	synthesizer := services.NewSynthesizer(wlParams)

	tmpl := &models.SQLTemplate{
		RawSQL:     "SELECT * FROM users WHERE id = :id AND region = :region",
		GroupKey:   "group1",
		Parameters: []string{":id", ":region"},
	}

	t.Run("Generate Single Query", func(t *testing.T) {
		sql, args, err := synthesizer.FillParameters(tmpl)
		assert.NoError(t, err)
		assert.Contains(t, sql, "SELECT * FROM users WHERE id = ")
		assert.Len(t, args, 2)
	})

	t.Run("Test Zipf Distribution Fidelity", func(t *testing.T) {
		counts := make(map[int]int)
		n := 1000

		for i := 0; i < n; i++ {
			_, args, err := synthesizer.FillParameters(tmpl)
			assert.NoError(t, err)
			idVal := args[0].(int) // id is first param
			counts[idVal]++
		}

		// With S=2.0, rank 0 (100) should be much more frequent than rank 1 (200)
		t.Logf("Counts: %v", counts)
		assert.Greater(t, counts[100], counts[200], "Rank 0 should be more frequent than Rank 1")
		assert.Greater(t, counts[200], counts[300], "Rank 1 should be more frequent than Rank 2 (or close to 0)")
	})
}

// TestGenerationFidelity simulates the E2E flow: Model -> Synthesizer -> Output Stats
func TestGenerationFidelity(t *testing.T) {
	// 1. Define Source Distribution (Zipf S=1.5)
	topValues := []interface{}{"A", "B", "C", "D", "E"}
	model := &models.ParameterModel{
		ParamName:   ":val",
		DistType:    models.DistZipfian,
		ZipfS:       1.5,
		Cardinality: 5,
		TopValues:   topValues,
	}

	wlParams := &models.WorkloadParameterModel{
		TemplateParameters: map[string]map[string]*models.ParameterModel{
			"g1": {":val": model},
		},
	}

	synth := services.NewSynthesizer(wlParams)
	tmpl := &models.SQLTemplate{GroupKey: "g1", Parameters: []string{":val"}, RawSQL: "SELECT :val"}

	// 2. Generate 1000 samples
	counts := make(map[string]int)
	n := 2000
	for i := 0; i < n; i++ {
		_, args, _ := synth.FillParameters(tmpl)
		val := args[0].(string)
		counts[val]++
	}

	// 3. Verify Rank Order (A > B > C > D > E)
	// Zipf(1.5, 5) probabilities roughly:
	// 1^-1.5 = 1
	// 2^-1.5 = 0.35
	// 3^-1.5 = 0.19
	// ...

	t.Logf("Generated Counts: %v", counts)

	assert.Greater(t, counts["A"], counts["B"])
	assert.Greater(t, counts["B"], counts["C"])

	// Basic Chi-Square Check or just ratio check
	ratio := float64(counts["A"]) / float64(counts["B"])
	expectedRatio := math.Pow(1.0, -1.5) / math.Pow(2.0, -1.5) // = 1 / 0.353 = 2.82

	// Allow some variance
	assert.InDelta(t, expectedRatio, ratio, 1.0, "Ratio between Rank 1 and 2 should be close to theoretical Zipf")
}
