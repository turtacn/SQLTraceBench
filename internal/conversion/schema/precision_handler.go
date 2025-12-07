package schema

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// PrecisionHandler handles precision logic for types like DECIMAL and TIMESTAMP.
type PrecisionHandler struct {
	policies map[string]PrecisionPolicy // Keyed by type category (e.g., "decimal", "timestamp")
    floatPolicy FloatPolicy
}

type FloatPolicy struct {
    PreferFloat64 bool `yaml:"prefer_float64"`
    WarnOnPrecisionLoss bool `yaml:"warn_on_precision_loss"`
    RoundingMode string `yaml:"rounding_mode"`
}

// NewPrecisionHandler creates a new PrecisionHandler.
func NewPrecisionHandler(policyPath string) *PrecisionHandler {
	h := &PrecisionHandler{
		policies: make(map[string]PrecisionPolicy),
	}
	h.loadPolicies(policyPath)
	return h
}

func (h *PrecisionHandler) loadPolicies(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		h.setDefaultPolicies()
		return
	}

	type Config struct {
		DecimalPolicy PrecisionPolicy `yaml:"decimal_policy"`
		TimestampPolicy struct {
			DefaultFractionalSeconds int `yaml:"default_fractional_seconds"`
            PreserveTimezone bool `yaml:"preserve_timezone"`
		} `yaml:"timestamp_policy"`
        FloatPolicy FloatPolicy `yaml:"float_policy"`
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		h.setDefaultPolicies()
		return
	}

	h.policies["decimal"] = cfg.DecimalPolicy
    h.floatPolicy = cfg.FloatPolicy
}

func (h *PrecisionHandler) setDefaultPolicies() {
	h.policies["decimal"] = PrecisionPolicy{
		MaxPrecision: 38,
		OverflowStrategy: "WARN",
	}
    h.floatPolicy = FloatPolicy{
        PreferFloat64: true,
        WarnOnPrecisionLoss: true,
    }
}

// HandlePrecision handles precision for a given type.
func (h *PrecisionHandler) HandlePrecision(baseType string, params []string, targetDB string) *PrecisionResult {
	result := &PrecisionResult{
		Warnings:    make([]TypeWarning, 0),
		Adjustments: make(map[string]any),
	}

	upperType := strings.ToUpper(baseType)

	if strings.Contains(upperType, "DECIMAL") || strings.Contains(upperType, "NUMERIC") {
		return h.handleDecimal(params, targetDB)
	} else if strings.Contains(upperType, "TIMESTAMP") || strings.Contains(upperType, "DATETIME") {
		return h.handleTimestamp(params, targetDB)
	} else if strings.Contains(upperType, "FLOAT") || strings.Contains(upperType, "DOUBLE") {
        // Only handle if we have a policy or targetDB specific logic, otherwise return type as is
        // But usually mapper handles this via direct mapping (DOUBLE -> Float64).
        // PrecisionHandler is for when parameters dictate the type.
        // For Float, parameters (M, D) are deprecated in MySQL but might exist.
        return h.handleFloat(baseType, params, targetDB)
    }

    result.TargetType = baseType
    if len(params) > 0 {
        result.TargetType = fmt.Sprintf("%s(%s)", baseType, strings.Join(params, ","))
    }
	return result
}

func (h *PrecisionHandler) handleDecimal(params []string, targetDB string) *PrecisionResult {
	result := &PrecisionResult{
		Warnings:    make([]TypeWarning, 0),
		Adjustments: make(map[string]any),
	}

	precision := 10 // default
	scale := 0

	if len(params) > 0 {
		precision, _ = strconv.Atoi(params[0])
	}
	if len(params) > 1 {
		scale, _ = strconv.Atoi(params[1])
	}

	policy := h.policies["decimal"]

	if precision > policy.MaxPrecision {
		result.HasLoss = true
		result.Warnings = append(result.Warnings, TypeWarning{
			Level:      "WARNING",
			Message:    fmt.Sprintf("Precision %d exceeds limit %d", precision, policy.MaxPrecision),
			Suggestion: "Consider using String for storage",
		})
        if precision <= 76 && targetDB == "clickhouse" {
             result.TargetType = fmt.Sprintf("Decimal256(%d)", scale)
             return result
        }
	}

    if targetDB == "clickhouse" {
        if precision <= 9 {
             result.TargetType = fmt.Sprintf("Decimal32(%d)", scale)
        } else if precision <= 18 {
             result.TargetType = fmt.Sprintf("Decimal64(%d)", scale)
        } else if precision <= 38 {
             result.TargetType = fmt.Sprintf("Decimal128(%d)", scale)
        } else {
             result.TargetType = fmt.Sprintf("Decimal256(%d)", scale)
        }
        return result
    }

	result.TargetType = fmt.Sprintf("Decimal(%d,%d)", precision, scale)
	return result
}

func (h *PrecisionHandler) handleTimestamp(params []string, targetDB string) *PrecisionResult {
    result := &PrecisionResult{
		Warnings:    make([]TypeWarning, 0),
		Adjustments: make(map[string]any),
	}

    frac := 0
    if len(params) > 0 {
        frac, _ = strconv.Atoi(params[0])
    }

    if targetDB == "clickhouse" {
        if frac == 0 {
            result.TargetType = "DateTime"
        } else if frac <= 9 {
             result.TargetType = fmt.Sprintf("DateTime64(%d)", frac)
             if frac > 6 {
                 result.Warnings = append(result.Warnings, TypeWarning{
                     Level: "INFO",
                     Message: "Timestamp precision > 6 (microseconds) detected",
                 })
             }
        } else {
            result.TargetType = "DateTime64(9)"
            result.HasLoss = true
             result.Warnings = append(result.Warnings, TypeWarning{
                 Level: "WARNING",
                 Message: fmt.Sprintf("Timestamp precision %d truncated to 9", frac),
             })
        }
    } else {
         result.TargetType = "TIMESTAMP"
    }

    return result
}

func (h *PrecisionHandler) handleFloat(baseType string, params []string, targetDB string) *PrecisionResult {
    result := &PrecisionResult{
		Warnings:    make([]TypeWarning, 0),
		Adjustments: make(map[string]any),
	}

    // MySQL FLOAT(p) where p is precision in bits. 0-24 -> FLOAT, 25-53 -> DOUBLE
    // Or FLOAT(M,D) (deprecated) -> Float32/64 depending on M.

    // ClickHouse: Float32, Float64

    isDouble := strings.Contains(strings.ToUpper(baseType), "DOUBLE")

    if targetDB == "clickhouse" {
        if isDouble || h.floatPolicy.PreferFloat64 {
            result.TargetType = "Float64"
        } else {
            result.TargetType = "Float32"
        }
    } else {
        result.TargetType = baseType // Pass through
    }
    return result
}
