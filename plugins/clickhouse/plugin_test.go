package clickhouse

import (
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {
	p := New()
	assert.NotNil(t, p)
	assert.Equal(t, "clickhouse", p.Name())
}