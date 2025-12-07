package schema

// TypeMappingContext holds context information for intelligent type mapping.
type TypeMappingContext struct {
	SourceType    string            // e.g. "VARCHAR(255)"
	SourceDB      string            // e.g. "mysql"
	TargetDB      string            // e.g. "clickhouse"
	ColumnName    string            // Column name for semantic analysis
	IsPrimaryKey  bool              // Is this a primary key column?
	IsIndexColumn bool              // Is this an index column?
	IsNullable    bool              // Is this column nullable?
	DefaultValue  string            // Default value
	TableContext  *TableContext     // Table level context
	CustomRules   map[string]string // Custom user rules
}

// TableContext holds table-level information.
type TableContext struct {
	TableName string
	Engine    string
	Charset   string
	Collation string
	// Add other table properties as needed
}

// TypeMappingResult holds the result of type mapping.
type TypeMappingResult struct {
	TargetType     string            // e.g. "String"
	Warnings       []TypeWarning     // Warnings generated during mapping
	Suggestions    []string          // Optimization suggestions
	PrecisionLoss  bool              // Indicates if precision loss occurred
	RequiresManual bool              // Indicates if manual review is required
	Metadata       map[string]any    // Additional metadata (e.g., original precision)
}

// TypeWarning represents a warning during type conversion.
type TypeWarning struct {
	Level          string // "INFO", "WARNING", "ERROR"
	Message        string
	Suggestion     string
	AffectedColumn string
}

// PrecisionPolicy defines policy for handling precision.
type PrecisionPolicy struct {
	MaxPrecision     int    `yaml:"max_precision"`
	MaxScale         int    `yaml:"max_scale"`
	OverflowStrategy string `yaml:"overflow_strategy"` // "TRUNCATE"/"ERROR"/"WARN"
	RoundingMode     string `yaml:"rounding_mode"`     // "ROUND_HALF_UP"/"FLOOR"/"CEIL"
}

// PrecisionResult holds the result of precision handling.
type PrecisionResult struct {
	TargetType  string
	HasLoss     bool
	Warnings    []TypeWarning
	Adjustments map[string]any
}

// AnalysisResult holds the result of type compatibility analysis.
type AnalysisResult struct {
	IsCompatible bool
	Warnings     []TypeWarning
	Suggestions  []string
	RiskLevel    string // "LOW", "MEDIUM", "HIGH"
}

// CompatibilityIssue describes a known compatibility issue.
type CompatibilityIssue struct {
	SourcePattern string
	TargetType    string
	IssueType     string // "OVERFLOW", "PRECISION_LOSS", "DATA_LOSS"
	Description   string
	Mitigation    string
}
