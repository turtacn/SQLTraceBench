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

	// We expect two unique values for the :id parameter
	assert.Len(t, paramDist.Values, 2, "should have 2 unique values")
	assert.Equal(t, 3, paramDist.Total, "should have 3 total observations")
}