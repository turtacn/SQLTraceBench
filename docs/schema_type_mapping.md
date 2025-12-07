# Schema Type Mapping

This document describes the type mapping rules for converting database schemas between different dialects.

## MySQL Source

| MySQL Type | ClickHouse Target | StarRocks Target | Notes |
|:-----------|:------------------|:-----------------|:------|
| TINYINT | Int8 | TINYINT | |
| SMALLINT | Int16 | SMALLINT | |
| INT / INTEGER | Int32 | INT | |
| BIGINT | Int64 | BIGINT | |
| FLOAT | Float32 | FLOAT | |
| DOUBLE | Float64 | DOUBLE | |
| DECIMAL(P,S) | Decimal(P,S) | DECIMAL(P,S) | |
| CHAR(N) | FixedString(N) | CHAR(N) | |
| VARCHAR | String | VARCHAR | |
| TEXT | String | STRING | |
| DATETIME | DateTime | DATETIME | |
| DATE | Date | DATE | |
| ENUM | Enum8 / Enum16 | VARCHAR | ClickHouse uses Enum8 if <= 127 items, else Enum16. StarRocks falls back to VARCHAR. |
| SET | String | VARCHAR | |
| BLOB | String | STRING | |
| BIT | UInt64 | BIGINT | |
| BOOLEAN | UInt8 | BOOLEAN | |

## PostgreSQL Source

| PostgreSQL Type | ClickHouse Target | StarRocks Target | Notes |
|:----------------|:------------------|:-----------------|:------|
| SMALLINT | Int16 | SMALLINT | |
| INTEGER | Int32 | INT | |
| BIGINT | Int64 | BIGINT | |
| REAL | Float32 | FLOAT | |
| DOUBLE PRECISION| Float64 | DOUBLE | |
| NUMERIC / DECIMAL| Decimal | DECIMAL | |
| VARCHAR | String | VARCHAR | |
| TEXT | String | STRING | |
| TIMESTAMP | DateTime | DATETIME | |
| DATE | Date | DATE | |
| BOOLEAN | UInt8 | BOOLEAN | |
| UUID | UUID | STRING | |
| JSON / JSONB | String | JSON | |
| ARRAY | Array(T) | - | StarRocks does not fully support nested types via simple mapping yet. |
| INET | IPv4 | VARCHAR | |
| MACADDR | String | VARCHAR | |
| SERIAL | Int32 | INT | Auto-increment removed |
| BIGSERIAL | Int64 | BIGINT | Auto-increment removed |

## TiDB Source

TiDB source largely follows MySQL mapping, with specific handling for:

* **Clustered Index**: Converted to `ORDER BY` in ClickHouse.
* **SHARD_ROW_ID_BITS**: Removed.
* **AUTO_RANDOM**: Removed.
