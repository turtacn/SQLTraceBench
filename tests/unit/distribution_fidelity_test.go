package unit

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
)

func TestZipfDistributionFidelity(t *testing.T) {
	// 1. Setup: Create a workload model with a parameter following a Zipfian distribution.
	paramName := "user_id"
	groupKey := "SELECT_FROM_USERS"
	numValues := 1000
	zipfS := 1.1

	// Create a parameter model with TopValues for the synthesizer to use.
	paramModel := &models.ParameterModel{
		ParamName:      paramName,
		DistType:       models.DistZipfian,
		ZipfS:          zipfS,
		TopValues:      make([]interface{}, numValues),
		TopFrequencies: make([]int, numValues), // Synthesizer doesn't use this, but good practice
	}
	for i := 0; i < numValues; i++ {
		paramModel.TopValues[i] = fmt.Sprintf("user-%d", i)
		// Frequencies aren't used by the sampler, but let's set them logically
		paramModel.TopFrequencies[i] = numValues - i
	}

	workloadModel := &models.WorkloadParameterModel{
		TemplateParameters: map[string]map[string]*models.ParameterModel{
			groupKey: {
				paramName: paramModel,
			},
		},
	}

	// Create a synthesizer instance.
	synthesizer := services.NewSynthesizer(workloadModel)

	// Create a dummy template that uses this parameter.
	template := &models.SQLTemplate{
		GroupKey:   groupKey,
		Parameters: []string{paramName},
	}

	// 2. Action: Generate a large number of parameter samples.
	numSamples := 20000
	frequencyMap := make(map[interface{}]int)

	for i := 0; i < numSamples; i++ {
		args, err := synthesizer.FillParameters(template)
		assert.NoError(t, err)
		assert.Len(t, args, 1)
		frequencyMap[args[0]]++
	}

	// 3. Verify: Check if the distribution of generated values approximates a Zipf distribution.

	// a. Check if the top-ranked item is the most frequent.
	topRankFreq := frequencyMap["user-0"]
	isTopRankMostFrequent := true
	for val, freq := range frequencyMap {
		if val != "user-0" && freq > topRankFreq {
			isTopRankMostFrequent = false
			break
		}
	}
	assert.True(t, isTopRankMostFrequent, "The top-ranked value ('user-0') should be the most frequent")

	// b. More rigorous check: Top 10% of items should account for a significant portion of samples.
	top10PercentCount := int(math.Ceil(0.1 * float64(numValues)))
	top10FrequencySum := 0
	for i := 0; i < top10PercentCount; i++ {
		val := fmt.Sprintf("user-%d", i)
		top10FrequencySum += frequencyMap[val]
	}

	proportion := float64(top10FrequencySum) / float64(numSamples)
	t.Logf("Top 10%% of items account for %.2f%% of samples", proportion*100)
	assert.True(t, proportion > 0.5, "Expected top 10%% of items to account for > 50%% of samples, but got %.2f%%", proportion*100)

	// c. Check if frequency is generally decreasing with rank.
	lastFreq := frequencyMap["user-0"]
	for i := 1; i < 20; i++ { // Check the top 20 ranks
		val := fmt.Sprintf("user-%d", i)
		currentFreq := frequencyMap[val]
		// Allow for some stochastic variance, but the trend should be downward.
		// A simple check is that the current frequency is not significantly larger than the last.
		assert.True(t, float64(currentFreq) < float64(lastFreq)*1.5, "Frequency should generally decrease with rank. Rank %d freq (%d) vs rank %d freq (%d)", i, currentFreq, i-1, lastFreq)
		lastFreq = currentFreq
	}
}
