package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TypeAnalyzer implements compatibility analysis.
type TypeAnalyzer struct {
	knownIssues map[string][]CompatibilityIssue
}

// NewTypeAnalyzer creates a new TypeAnalyzer.
func NewTypeAnalyzer() *TypeAnalyzer {
	analyzer := &TypeAnalyzer{
		knownIssues: make(map[string][]CompatibilityIssue),
	}
	analyzer.initKnownIssues()
	return analyzer
}

func (a *TypeAnalyzer) initKnownIssues() {
	a.knownIssues = map[string][]CompatibilityIssue{
		"mysql:clickhouse": {
			{
				SourcePattern: "BIGINT UNSIGNED",
				TargetType:    "Int64",
				IssueType:     "OVERFLOW",
				Description:   "Values > 2^63-1 will cause overflow in Int64",
				Mitigation:    "Use UInt64",
			},
			{
				SourcePattern: "DOUBLE",
				TargetType:    "Float32",
				IssueType:     "PRECISION_LOSS",
				Description:   "Float32 has lower precision than DOUBLE",
				Mitigation:    "Use Float64",
			},
			{
				SourcePattern: "GEOMETRY",
				TargetType:    "String",
				IssueType:     "DATA_LOSS",
				Description:   "Spatial data stored as text, indexes lost",
				Mitigation:    "Use specialized Geo types if available or store as WKT",
			},
		},
	}
}

// Analyze performs compatibility analysis.
func (a *TypeAnalyzer) Analyze(sourceType, targetType string, ctx *TypeMappingContext) *AnalysisResult {
	result := &AnalysisResult{
		IsCompatible: true,
		Warnings:     make([]TypeWarning, 0),
		Suggestions:  make([]string, 0),
		RiskLevel:    "LOW",
	}

	// 1. Check Known Issues
	key := fmt.Sprintf("%s:%s", ctx.SourceDB, ctx.TargetDB)
	if issues, ok := a.knownIssues[key]; ok {
		for _, issue := range issues {
			if strings.Contains(strings.ToUpper(sourceType), issue.SourcePattern) &&
               strings.Contains(targetType, issue.TargetType) {
				result.Warnings = append(result.Warnings, TypeWarning{
					Level:          "WARNING",
					Message:        issue.Description,
					Suggestion:     issue.Mitigation,
					AffectedColumn: ctx.ColumnName,
				})
                if issue.IssueType == "OVERFLOW" || issue.IssueType == "DATA_LOSS" {
                    result.RiskLevel = "HIGH"
                } else {
                    result.RiskLevel = "MEDIUM"
                }
			}
		}
	}

    // 2. Numeric Overflow Analysis
    if isNumericType(sourceType) && isNumericType(targetType) {
        if canOverflow(sourceType, targetType) {
            result.Warnings = append(result.Warnings, TypeWarning{
                Level: "WARNING",
                Message: fmt.Sprintf("%s may overflow when converting to %s", sourceType, targetType),
                Suggestion: "Add data validation or use larger target type",
            })
            if result.RiskLevel == "LOW" {
                result.RiskLevel = "MEDIUM"
            }
        }
    }

	// 3. Precision Loss Analysis
    if (isDecimalType(sourceType) || isFloatType(sourceType)) && (isDecimalType(targetType) || isFloatType(targetType)) {
        sP, sS := extractDecimalParams(sourceType)
        tP, tS := extractDecimalParams(targetType)

        // Adjust for ClickHouse implicit precision
        if tP == 0 && tS > 0 {
             if strings.Contains(targetType, "Decimal32") {
                 tP = 9
                 tS = tPFromParam(targetType)
             } else if strings.Contains(targetType, "Decimal64") {
                 tP = 18
                 tS = tPFromParam(targetType)
             } else if strings.Contains(targetType, "Decimal128") {
                 tP = 38
                 tS = tPFromParam(targetType)
             } else if strings.Contains(targetType, "Decimal256") {
                 tP = 76
                 tS = tPFromParam(targetType)
             }
        }

        if sP > 0 && tP > 0 && sP > tP {
             result.Warnings = append(result.Warnings, TypeWarning{
                Level: "WARNING",
                Message: fmt.Sprintf("Precision loss: %s(%d) -> %s(%d)", sourceType, sP, targetType, tP),
                Suggestion: "Consider using higher precision target type",
            })
             result.RiskLevel = "HIGH"
        }

        if sS > 0 && tS > 0 && sS > tS {
             result.Warnings = append(result.Warnings, TypeWarning{
                Level: "WARNING",
                Message: fmt.Sprintf("Scale loss: %s(...,%d) -> %s(...,%d)", sourceType, sS, targetType, tS),
                Suggestion: "Consider using higher scale target type",
            })
             result.RiskLevel = "HIGH"
        }
    }

	// 4. Timezone Analysis
    if isTimestampType(sourceType) && strings.Contains(strings.ToUpper(sourceType), "TZ") && !strings.Contains(targetType, "TZ") && !strings.Contains(targetType, "DateTime64") {
          result.Warnings = append(result.Warnings, TypeWarning{
            Level: "INFO",
            Message: "Timestamp with timezone detected",
            Suggestion: "Use DateTime('UTC') or DateTime64 to preserve timezone info better",
        })
    }

	return result
}

func tPFromParam(t string) int {
    start := strings.Index(t, "(")
    end := strings.LastIndex(t, ")")
    if start != -1 && end != -1 {
        val, _ := strconv.Atoi(t[start+1 : end])
        return val
    }
    return 0
}

func isNumericType(t string) bool {
    u := strings.ToUpper(t)
    return strings.Contains(u, "INT") || strings.Contains(u, "FLOAT") || strings.Contains(u, "DOUBLE") || strings.Contains(u, "DECIMAL") || strings.Contains(u, "NUMERIC")
}

func isDecimalType(t string) bool {
    return strings.Contains(strings.ToUpper(t), "DECIMAL") || strings.Contains(strings.ToUpper(t), "NUMERIC")
}

func isFloatType(t string) bool {
    return strings.Contains(strings.ToUpper(t), "FLOAT") || strings.Contains(strings.ToUpper(t), "DOUBLE") || strings.Contains(strings.ToUpper(t), "REAL")
}

func isTimestampType(t string) bool {
    return strings.Contains(strings.ToUpper(t), "TIMESTAMP") || strings.Contains(strings.ToUpper(t), "DATETIME")
}

func canOverflow(source, target string) bool {
    s := strings.ToUpper(source)
    t := strings.ToUpper(target)
    if strings.Contains(s, "UNSIGNED") && strings.Contains(s, "BIGINT") && t == "INT64" {
        return true
    }
    return false
}

func extractDecimalParams(t string) (int, int) {
    base, params := analyzerParseTypeWithParams(t)

    // Check specific CH types first
    if strings.Contains(base, "Decimal32") { return 9, getScale(params) }
    if strings.Contains(base, "Decimal64") { return 18, getScale(params) }
    if strings.Contains(base, "Decimal128") { return 38, getScale(params) }
    if strings.Contains(base, "Decimal256") { return 76, getScale(params) }

    if len(params) >= 2 {
        p, _ := strconv.Atoi(params[0])
        s, _ := strconv.Atoi(params[1])
        return p, s
    } else if len(params) == 1 {
        p, _ := strconv.Atoi(params[0])
        return p, 0
    }
    return 0, 0
}

func getScale(params []string) int {
    if len(params) > 0 {
        v, _ := strconv.Atoi(params[0])
        return v
    }
    return 0
}

func analyzerParseTypeWithParams(fullType string) (string, []string) {
	re := regexp.MustCompile(`^([a-zA-Z0-9_ ]+)(?:\(([^)]+)\))?.*$`)
	matches := re.FindStringSubmatch(fullType)
	if len(matches) < 2 {
		return fullType, nil
	}
	baseType := strings.TrimSpace(matches[1])
	var params []string
	if len(matches) > 2 && matches[2] != "" {
		rawParams := strings.Split(matches[2], ",")
		for _, p := range rawParams {
			params = append(params, strings.TrimSpace(p))
		}
	}
	return baseType, params
}
