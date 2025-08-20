package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/turtacn/SQLTraceBench/pkg/types"
)

type SQLTemplate struct {
	TemplateID    string          `json:"template_id"`
	OriginalQuery string          `json:"original_query"`
	TemplateQuery string          `json:"template_query"`
	Parameters    []Parameter     `json:"parameters"`
	QueryType     types.QueryType `json:"query_type"`
	Tables        []string        `json:"tables"`
	Frequency     int64           `json:"frequency"`
	CreatedAt     time.Time       `json:"created_at"`
}

type Parameter struct {
	Name     string              `json:"name"`
	Type     types.ParameterType `json:"type"`
	Position int                 `json:"position"`
}

func (t *SQLTemplate) Validate() *types.SQLTraceBenchError {
	if t.OriginalQuery == "" || t.TemplateQuery == "" {
		return types.NewError(types.ErrInvalidInput, "query fields cannot be empty")
	}
	return nil
}

func (t *SQLTemplate) GenerateQuery(params map[string]interface{}) (string, error) {
	query := t.TemplateQuery
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	query = re.ReplaceAllStringFunc(query, func(s string) string {
		name := strings.Trim(s, "{}")
		if val, ok := params[name]; ok {
			return fmt.Sprintf("'%v'", val)
		}
		return s
	})
	return query, nil
}

func (t *SQLTemplate) ExtractParameters() map[string]interface{} {
	return nil // TODO: actual extraction logic
}

//Personal.AI order the ending
