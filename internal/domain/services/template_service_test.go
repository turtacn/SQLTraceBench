package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/SQLTraceBench/internal/domain/models"
)

func TestTemplateService_ExtractTemplates(t *testing.T) {
	service := NewTemplateService()

	// Create a collection of traces with duplicates and variations
	tc := models.TraceCollection{}
	tc.Add(models.SQLTrace{Query: "select * from users where id = :id"})
	tc.Add(models.SQLTrace{Query: "SELECT * FROM users WHERE id = :id "}) // a duplicate with different case/spacing
	tc.Add(models.SQLTrace{Query: "select * from orders"})
	tc.Add(models.SQLTrace{Query: "select * from users where id = :id"}) // another duplicate

	templates := service.ExtractTemplates(tc)

	// We expect two unique templates
	assert.Len(t, templates, 2, "should have 2 unique templates")

	// The first template should be the most frequent one
	usersTemplate := templates[0]
	assert.Equal(t, 3, usersTemplate.Weight, "users template should have weight 3")
	assert.Equal(t, "select * from users where id = :id", usersTemplate.GroupKey, "group key should be normalized")
	assert.Contains(t, usersTemplate.Parameters, ":id", "should have extracted :id parameter")

	// The second template should be the less frequent one
	ordersTemplate := templates[1]
	assert.Equal(t, 1, ordersTemplate.Weight, "orders template should have weight 1")
	assert.Equal(t, "select * from orders", ordersTemplate.GroupKey, "group key should be normalized")
	assert.Empty(t, ordersTemplate.Parameters, "should have no parameters")
}