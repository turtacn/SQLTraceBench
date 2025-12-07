package schema

import "strings"

// TypeMapping defines the mapping from source types to target types.
type TypeMapping map[string]string

var (
	// MySQL to ClickHouse mapping
	mysqlToClickHouseTypeMap = TypeMapping{
		"TINYINT":    "Int8",
		"SMALLINT":   "Int16",
		"MEDIUMINT":  "Int32",
		"INT":        "Int32",
		"INTEGER":    "Int32",
		"BIGINT":     "Int64",
		"FLOAT":      "Float32",
		"DOUBLE":     "Float64",
		"DECIMAL":    "Decimal", // Requires P, S
		"CHAR":       "FixedString", // Requires N
		"VARCHAR":    "String",
		"TEXT":       "String",
		"TINYTEXT":   "String",
		"MEDIUMTEXT": "String",
		"LONGTEXT":   "String",
		"DATETIME":   "DateTime",
		"TIMESTAMP":  "DateTime",
		"DATE":       "Date",
		"JSON":       "String",
		"ENUM":       "Enum8",
		"SET":        "String",
		"BLOB":       "String",
		"TINYBLOB":   "String",
		"MEDIUMBLOB": "String",
		"LONGBLOB":   "String",
		"BIT":        "UInt64", // Simplified
		"BOOLEAN":    "UInt8",
	}

	// MySQL to StarRocks mapping
	mysqlToStarRocksTypeMap = TypeMapping{
		"TINYINT":    "TINYINT",
		"SMALLINT":   "SMALLINT",
		"MEDIUMINT":  "INT",
		"INT":        "INT",
		"INTEGER":    "INT",
		"BIGINT":     "BIGINT",
		"FLOAT":      "FLOAT",
		"DOUBLE":     "DOUBLE",
		"DECIMAL":    "DECIMAL",
		"CHAR":       "CHAR",
		"VARCHAR":    "VARCHAR",
		"TEXT":       "STRING",
		"TINYTEXT":   "STRING",
		"MEDIUMTEXT": "STRING",
		"LONGTEXT":   "STRING",
		"DATETIME":   "DATETIME",
		"TIMESTAMP":  "DATETIME",
		"DATE":       "DATE",
		"JSON":       "JSON",
		"ENUM":       "VARCHAR", // StarRocks doesn't fully support ENUM yet, usually VARCHAR
		"SET":        "VARCHAR",
		"BLOB":       "STRING",
		"BIT":        "BIGINT",
		"BOOLEAN":    "BOOLEAN",
	}

	// Postgres to ClickHouse mapping
	postgresToClickHouseTypeMap = TypeMapping{
		"SMALLINT":         "Int16",
		"INTEGER":          "Int32",
		"INT":              "Int32",
		"BIGINT":           "Int64",
		"REAL":             "Float32",
		"DOUBLE PRECISION": "Float64",
		"NUMERIC":          "Decimal",
		"DECIMAL":          "Decimal",
		"VARCHAR":          "String",
		"CHARACTER VARYING": "String",
		"CHAR":             "FixedString",
		"CHARACTER":        "FixedString",
		"TEXT":             "String",
		"TIMESTAMP":        "DateTime",
		"TIMESTAMPTZ":      "DateTime", // DateTime('UTC') handled in code
		"DATE":             "Date",
		"BOOLEAN":          "UInt8",
		"BOOL":             "UInt8",
		"UUID":             "UUID",
		"JSON":             "String",
		"JSONB":            "String",
		"INET":             "IPv4", // or IPv6, requires logic
		"MACADDR":          "String",
		"SERIAL":           "Int32",
		"BIGSERIAL":        "Int64",
	}

	// Postgres to StarRocks mapping
	postgresToStarRocksTypeMap = TypeMapping{
		"SMALLINT":         "SMALLINT",
		"INTEGER":          "INT",
		"INT":              "INT",
		"BIGINT":           "BIGINT",
		"REAL":             "FLOAT",
		"DOUBLE PRECISION": "DOUBLE",
		"NUMERIC":          "DECIMAL",
		"DECIMAL":          "DECIMAL",
		"VARCHAR":          "VARCHAR",
		"TEXT":             "STRING",
		"TIMESTAMP":        "DATETIME",
		"TIMESTAMPTZ":      "DATETIME",
		"DATE":             "DATE",
		"BOOLEAN":          "BOOLEAN",
		"BOOL":             "BOOLEAN",
		"UUID":             "STRING", // StarRocks doesn't have native UUID
		"JSON":             "JSON",
		"JSONB":            "JSON",
		"INET":             "VARCHAR",
		"MACADDR":          "VARCHAR",
		"SERIAL":           "INT",
		"BIGSERIAL":        "BIGINT",
	}
)

func getBaseType(fullType string) string {
	idx := strings.Index(fullType, "(")
	if idx != -1 {
		return strings.ToUpper(strings.TrimSpace(fullType[:idx]))
	}
	return strings.ToUpper(strings.TrimSpace(fullType))
}
