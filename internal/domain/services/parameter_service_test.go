package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestParameterService_BuildModel(t *testing.T) {
	service := NewParameterService()

	// Create a collection of traces and templates
	tc := models.TraceCollection{}
	tc.Add(models.SQLTrace{Query: "select * from users where id = 1"})
	tc.Add(models.SQLTrace{Query: "select * from users where id = 2"})
	tc.Add(models.SQLTrace{Query: "select * from users where id = 1"})

	// Ensure we match what the extractor expects.
	// The extractor expects :id, but trace uses 1.
	// Regex extractor might expect params in raw SQL or similar.
	// But in this test, we construct SQLTemplate manually.

	templates := []models.SQLTemplate{
		{
			RawSQL:   "select * from users where id = :id",
			GroupKey: "select * from users where id = :id",
			Parameters: []string{":id"},
		},
	}

	pm := service.BuildModel(tc, templates)

	// We expect one template in the model
	assert.Len(t, pm.TemplateParameters, 1, "should have 1 template in the model")

	// Get the parameter distribution for the :id parameter
	paramDist, ok := pm.TemplateParameters["select * from users where id = :id"][":id"]
	assert.True(t, ok, "should have a distribution for the :id parameter")

	// We expect two unique values for the :id parameter (1 and 2)
	// But note: "1" appears twice, "2" once.
	// HotspotDetector puts them in TopValues.
	assert.Len(t, paramDist.TopValues, 2, "should have 2 unique top values")

	// Verify counts if possible via TopFrequencies
	// "1" count 2, "2" count 1.
	// TopValues should be sorted by frequency.
	assert.Equal(t, "1", paramDist.TopValues[0])
	assert.Equal(t, 2, paramDist.TopFrequencies[0])
}
