package types

import "strings"

type DatabaseType int

const (
	DatabaseNone DatabaseType = iota
	DatabaseStarRocks
	DatabaseClickHouse
	DatabaseMySQL
	DatabasePostgreSQL
	DatabaseTiDB
	DatabaseDoris
	DatabaseMongoDB
)

func (d DatabaseType) String() string {
	return [...]string{"", "starrocks", "clickhouse", "mysql", "postgresql", "tidb", "doris", "mongodb"}[d]
}

func DatabaseTypeFromString(s string) DatabaseType {
	switch strings.ToLower(s) {
	case "starrocks":
		return DatabaseStarRocks
	case "clickhouse":
		return DatabaseClickHouse
	case "mysql":
		return DatabaseMySQL
	case "postgresql", "postgres":
		return DatabasePostgreSQL
	case "tidb":
		return DatabaseTiDB
	case "doris":
		return DatabaseDoris
	case "mongodb", "mongo":
		return DatabaseMongoDB
	default:
		return DatabaseNone
	}
}

type QueryType int

const (
	QuerySelect QueryType = iota
	QueryInsert
	QueryUpdate
	QueryDelete
	QueryDDL
	QueryOther
)

func (q QueryType) String() string {
	return [...]string{"SELECT", "INSERT", "UPDATE", "DELETE", "DDL", "OTHER"}[q]
}

type ParameterType int

const (
	TypeInteger ParameterType = iota
	TypeString
	TypeFloat
	TypeBoolean
	TypeDateTime
	TypeJSON
)

func (p ParameterType) String() string {
	return [...]string{"int", "string", "float", "bool", "datetime", "json"}[p]
}

type DistributionType int

const (
	DistributionUniform DistributionType = iota
	DistributionNormal
	DistributionZipfian
	DistributionExponential
)

func (d DistributionType) String() string {
	return [...]string{"uniform", "normal", "zipfian", "exponential"}[d]
}

type BenchmarkStatus int

const (
	StatusPending BenchmarkStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusCancelled
)

func (b BenchmarkStatus) String() string {
	return [...]string{"pending", "running", "completed", "failed", "cancelled"}[b]
}

//Personal.AI order the ending
