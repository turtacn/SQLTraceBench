package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// IntelligentTypeMapper implements context-aware type mapping.
type IntelligentTypeMapper struct {
    ruleLoader *MappingRuleLoader
	analyzer    *TypeAnalyzer
	precision   *PrecisionHandler
	mutex       sync.RWMutex
}

// NewIntelligentTypeMapper creates a new IntelligentTypeMapper.
func NewIntelligentTypeMapper(
    ruleLoader *MappingRuleLoader,
	analyzer *TypeAnalyzer,
	precision *PrecisionHandler,
) *IntelligentTypeMapper {
	return &IntelligentTypeMapper{
        ruleLoader:  ruleLoader,
		analyzer:    analyzer,
		precision:   precision,
	}
}

// MapType performs intelligent type mapping.
func (m *IntelligentTypeMapper) MapType(ctx *TypeMappingContext) (*TypeMappingResult, error) {
	result := &TypeMappingResult{
		Warnings:    make([]TypeWarning, 0),
		Suggestions: make([]string, 0),
		Metadata:    make(map[string]any),
	}

	baseType, params := mapperParseTypeWithParams(ctx.SourceType)

	// 1. Check Custom Rules (Highest Priority)
	if customMapping := m.getCustomMapping(ctx.SourceDB, baseType, ctx.TargetDB); customMapping != "" {
		result.TargetType = applyParams(customMapping, params)
		result.Suggestions = append(result.Suggestions, "Using custom mapping rule")
		return result, nil
	}

	// 2. Context-Aware Mapping (Dynamic from RuleLoader)
    if match := m.ruleLoader.MatchContextRules(ctx); match != "" {
        result.TargetType = match
        result.Suggestions = append(result.Suggestions, fmt.Sprintf("Applied context rule"))
        return result, nil
    }

	// 3. Base Type Mapping
	baseMapping := m.getBaseMapping(ctx.SourceDB, baseType, ctx.TargetDB)
	if baseMapping == "" {
		result.TargetType = "String"
		result.Warnings = append(result.Warnings, TypeWarning{
			Level:      "ERROR",
			Message:    fmt.Sprintf("Unknown type '%s', fallback to String", baseType),
			Suggestion: "Add custom mapping rule",
		})
		result.RequiresManual = true
		return result, nil
	}

	// 4. Precision Handling
	if needsPrecisionHandling(baseType) {
		precisionResult := m.precision.HandlePrecision(baseType, params, ctx.TargetDB)
		result.TargetType = precisionResult.TargetType
		result.PrecisionLoss = precisionResult.HasLoss
		result.Warnings = append(result.Warnings, precisionResult.Warnings...)
		result.Metadata["original_params"] = params
	} else {
		result.TargetType = applyParams(baseMapping, params)
	}

	// 5. Compatibility Analysis
	analysisResult := m.analyzer.Analyze(ctx.SourceType, result.TargetType, ctx)
	result.Warnings = append(result.Warnings, analysisResult.Warnings...)
	result.Suggestions = append(result.Suggestions, analysisResult.Suggestions...)

	return result, nil
}

func (m *IntelligentTypeMapper) getBaseMapping(sourceDB, sourceType, targetDB string) string {
    rules := m.ruleLoader.GetRules()
    if rules == nil {
        return ""
    }

	key := fmt.Sprintf("%s:%s", sourceDB, targetDB)
	if mappings, ok := rules.DefaultRules[key]; ok {
		if mapping, found := mappings[strings.ToUpper(sourceType)]; found {
            return mapping
        }
        for k, v := range mappings {
            if strings.EqualFold(k, sourceType) {
                return v
            }
        }
	}
	return ""
}

func (m *IntelligentTypeMapper) getCustomMapping(sourceDB, sourceType, targetDB string) string {
    rules := m.ruleLoader.GetRules()
    if rules == nil {
        return ""
    }

	key := fmt.Sprintf("%s:%s", sourceDB, targetDB)
	if mappings, ok := rules.CustomRules[key]; ok {
		if mapping, found := mappings[strings.ToUpper(sourceType)]; found {
            return mapping
        }
        for k, v := range mappings {
            if strings.EqualFold(k, sourceType) {
                return v
            }
        }
	}
	return ""
}

func mapperParseTypeWithParams(fullType string) (string, []string) {
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

func applyParams(targetType string, params []string) string {
    if targetType == "String" {
        return targetType
    }
    if len(params) == 0 {
        return targetType
    }
    if !strings.Contains(targetType, "(") {
        return fmt.Sprintf("%s(%s)", targetType, strings.Join(params, ","))
    }
    return targetType
}

func extractLength(params []string) int {
	if len(params) > 0 {
		if val, err := strconv.Atoi(params[0]); err == nil {
			return val
		}
	}
	return 0
}

func needsPrecisionHandling(baseType string) bool {
    upper := strings.ToUpper(baseType)
    return strings.Contains(upper, "DECIMAL") ||
           strings.Contains(upper, "TIMESTAMP") ||
           strings.Contains(upper, "DATETIME") ||
           strings.Contains(upper, "FLOAT") ||
           strings.Contains(upper, "DOUBLE")
}
