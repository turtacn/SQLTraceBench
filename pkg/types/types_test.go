package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseType(t *testing.T) {
	assert.Equal(t, "starrocks", DatabaseStarRocks.String())
	assert.Equal(t, "clickhouse", DatabaseClickHouse.String())
	assert.Equal(t, "mysql", DatabaseMySQL.String())
	assert.Equal(t, "postgresql", DatabasePostgreSQL.String())
	assert.Equal(t, "tidb", DatabaseTiDB.String())
	assert.Equal(t, "doris", DatabaseDoris.String())
	assert.Equal(t, "mongodb", DatabaseMongoDB.String())

	assert.Equal(t, DatabaseStarRocks, DatabaseTypeFromString("starrocks"))
	assert.Equal(t, DatabaseClickHouse, DatabaseTypeFromString("clickhouse"))
	assert.Equal(t, DatabaseMySQL, DatabaseTypeFromString("mysql"))
	assert.Equal(t, DatabasePostgreSQL, DatabaseTypeFromString("postgresql"))
	assert.Equal(t, DatabasePostgreSQL, DatabaseTypeFromString("postgres"))
	assert.Equal(t, DatabaseTiDB, DatabaseTypeFromString("tidb"))
	assert.Equal(t, DatabaseDoris, DatabaseTypeFromString("doris"))
	assert.Equal(t, DatabaseMongoDB, DatabaseTypeFromString("mongodb"))
	assert.Equal(t, DatabaseMongoDB, DatabaseTypeFromString("mongo"))
	assert.Equal(t, DatabaseNone, DatabaseTypeFromString("unknown"))
}

func TestQueryType(t *testing.T) {
	assert.Equal(t, "SELECT", QuerySelect.String())
	assert.Equal(t, "INSERT", QueryInsert.String())
	assert.Equal(t, "UPDATE", QueryUpdate.String())
	assert.Equal(t, "DELETE", QueryDelete.String())
	assert.Equal(t, "DDL", QueryDDL.String())
	assert.Equal(t, "OTHER", QueryOther.String())
}

func TestParameterType(t *testing.T) {
	assert.Equal(t, "int", TypeInteger.String())
	assert.Equal(t, "string", TypeString.String())
	assert.Equal(t, "float", TypeFloat.String())
	assert.Equal(t, "bool", TypeBoolean.String())
	assert.Equal(t, "datetime", TypeDateTime.String())
	assert.Equal(t, "json", TypeJSON.String())
}

func TestDistributionType(t *testing.T) {
	assert.Equal(t, "uniform", DistributionUniform.String())
	assert.Equal(t, "normal", DistributionNormal.String())
	assert.Equal(t, "zipfian", DistributionZipfian.String())
	assert.Equal(t, "exponential", DistributionExponential.String())
}

func TestBenchmarkStatus(t *testing.T) {
	assert.Equal(t, "pending", StatusPending.String())
	assert.Equal(t, "running", StatusRunning.String())
	assert.Equal(t, "completed", StatusCompleted.String())
	assert.Equal(t, "failed", StatusFailed.String())
	assert.Equal(t, "cancelled", StatusCancelled.String())
}

func TestSQLTraceBenchError(t *testing.T) {
	// Test NewError
	err := NewError(ErrInvalidInput, "invalid input")
	assert.Equal(t, "invalid_input: invalid input", err.Error())

	// Test WrapError
	cause := errors.New("cause")
	err = WrapError(ErrInternal, "internal error", cause)
	assert.Equal(t, "internal: internal error (caused by: cause)", err.Error())
}