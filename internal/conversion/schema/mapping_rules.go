package schema

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/turtacn/SQLTraceBench/pkg/utils"
	"gopkg.in/yaml.v3"
)

// MappingRules defines the structure of the rules file.
type MappingRules struct {
	Version      string                       `yaml:"version"`
	UpdatedAt    time.Time                    `yaml:"updated_at"`
	DefaultRules map[string]map[string]string `yaml:"default_rules"`
	CustomRules  map[string]map[string]string `yaml:"custom_rules"`
	ContextRules []ContextMappingRule         `yaml:"context_rules"`
}

type ContextMappingRule struct {
	Name        string      `yaml:"name"`
	Conditions  []Condition `yaml:"conditions"`
	TargetType  string      `yaml:"target_type"`
	Priority    int         `yaml:"priority"`
	Description string      `yaml:"description"`
}

type Condition struct {
	Field    string `yaml:"field"` // "column_name"/"is_primary_key"/"source_type"/"length"
	Operator string `yaml:"operator"` // "equals"/"contains"/"matches"/"range"
	Value    any    `yaml:"value"`
}

// MappingRuleLoader loads and manages mapping rules.
type MappingRuleLoader struct {
	rulesPath   string
	rules       *MappingRules
	watcher     *fsnotify.Watcher
	reloadChan  chan bool
	mutex       sync.RWMutex
}

// NewMappingRuleLoader creates a new loader.
func NewMappingRuleLoader(rulesPath string) (*MappingRuleLoader, error) {
	loader := &MappingRuleLoader{
		rulesPath:  rulesPath,
		reloadChan: make(chan bool, 1),
	}

	// Try load, if fails, load defaults
	if err := loader.Load(); err != nil {
		utils.GetGlobalLogger().Warn("Failed to load rules file, using defaults", utils.Field{Key: "error", Value: err})
		loader.loadDefaults()
	}

	// Initialize watcher
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		loader.watcher = watcher
		if err := watcher.Add(rulesPath); err == nil {
			go loader.watch()
		} else {
            // It's okay if file doesn't exist yet, but we can't watch it.
            // If it exists later, we miss it unless we watch dir.
            // For simplicity, we ignore watcher error if file not found.
        }
	}

	return loader, nil
}

func (l *MappingRuleLoader) watch() {
	for {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
                // Debounce could be added here
				time.Sleep(100 * time.Millisecond)
				if err := l.Load(); err == nil {
					select {
					case l.reloadChan <- true:
					default:
					}
				}
			}
		case _, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

// Load loads the rules from file.
func (l *MappingRuleLoader) Load() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	data, err := os.ReadFile(l.rulesPath)
	if err != nil {
		return fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules MappingRules
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("failed to parse rules YAML: %w", err)
	}

	l.rules = &rules
	return nil
}

func (l *MappingRuleLoader) loadDefaults() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.rules = &MappingRules{
		Version: "1.0.0 (Default)",
		DefaultRules: map[string]map[string]string{
			"mysql:clickhouse": {
				"VARCHAR":   "String",
				"CHAR":      "FixedString",
				"INT":       "Int32",
				"INTEGER":   "Int32",
				"TINYINT":   "Int8",
				"SMALLINT":  "Int16",
				"MEDIUMINT": "Int32",
				"BIGINT":    "Int64",
				"DECIMAL":   "Decimal128",
				"TIMESTAMP": "DateTime",
				"DATETIME":  "DateTime",
				"DATE":      "Date",
				"TEXT":      "String",
				"BLOB":      "String",
				"ENUM":      "Enum8",
				"JSON":      "String",
				"FLOAT":     "Float32",
				"DOUBLE":    "Float64",
			},
			"postgres:clickhouse": {
				"VARCHAR":     "String",
				"INTEGER":     "Int32",
				"INT":         "Int32",
				"SMALLINT":    "Int16",
				"BIGINT":      "Int64",
				"NUMERIC":     "Decimal128",
				"DECIMAL":     "Decimal128",
				"TIMESTAMP":   "DateTime64",
				"TIMESTAMPTZ": "DateTime64",
				"DATE":        "Date",
				"TEXT":        "String",
				"JSONB":       "String",
				"JSON":        "String",
				"SERIAL":      "Int32",
			},
		},
		CustomRules: make(map[string]map[string]string),
        ContextRules: []ContextMappingRule{
            // Default context rules for testing/safety
            {
                Name: "primary_key_varchar_to_fixedstring",
                Priority: 100,
                Conditions: []Condition{
                    {Field: "is_primary_key", Operator: "equals", Value: true},
                    {Field: "source_type", Operator: "matches", Value: `VARCHAR\((\d+)\)`},
                    {Field: "length", Operator: "range", Value: []any{1, 256}},
                },
                TargetType: "FixedString({length})",
            },
        },
	}
}

// GetRules returns the current rules.
func (l *MappingRuleLoader) GetRules() *MappingRules {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.rules
}

// Subscribe returns a channel for updates.
func (l *MappingRuleLoader) Subscribe() <-chan bool {
	return l.reloadChan
}

// MatchContextRules finds the best matching rule for the context.
func (l *MappingRuleLoader) MatchContextRules(ctx *TypeMappingContext) string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

    if l.rules == nil {
        return ""
    }

	// Filter applicable rules
	var candidates []ContextMappingRule
	for _, rule := range l.rules.ContextRules {
		if l.matchesRule(rule, ctx) {
			candidates = append(candidates, rule)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	// Sort by priority desc
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Priority > candidates[j].Priority
	})

	// Use best match
	best := candidates[0]
	return l.substituteVariables(best.TargetType, ctx)
}

func (l *MappingRuleLoader) matchesRule(rule ContextMappingRule, ctx *TypeMappingContext) bool {
	for _, cond := range rule.Conditions {
		if !l.matchesCondition(cond, ctx) {
			return false
		}
	}
	return true
}

func (l *MappingRuleLoader) matchesCondition(cond Condition, ctx *TypeMappingContext) bool {
	var fieldValue any
	switch cond.Field {
	case "column_name":
		fieldValue = ctx.ColumnName
	case "source_type":
		fieldValue = ctx.SourceType
	case "is_primary_key":
		fieldValue = ctx.IsPrimaryKey
	case "is_nullable":
		fieldValue = ctx.IsNullable
	case "length":
        _, params := parseTypeWithParams(ctx.SourceType)
        if len(params) > 0 {
            if v, err := strconv.Atoi(params[0]); err == nil {
                fieldValue = v
            }
        }
	default:
		return false
	}

    if fieldValue == nil {
        return false
    }

	switch cond.Operator {
	case "equals":
		return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", cond.Value)
	case "contains":
		s, ok1 := fieldValue.(string)
		sub, ok2 := cond.Value.(string)
		if ok1 && ok2 {
			return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
		}
	case "matches":
		s, ok1 := fieldValue.(string)
		pat, ok2 := cond.Value.(string)
		if ok1 && ok2 {
			matched, _ := regexp.MatchString(pat, s)
			return matched
		}
    case "range":
        v, ok1 := fieldValue.(int)
        r, ok2 := cond.Value.([]any)
        if ok1 && ok2 && len(r) == 2 {
            min, _ := toInt(r[0])
            max, _ := toInt(r[1])
            return v >= min && v <= max
        }
	}
	return false
}

func (l *MappingRuleLoader) substituteVariables(targetType string, ctx *TypeMappingContext) string {
    res := targetType
    if strings.Contains(res, "{length}") {
        _, params := parseTypeWithParams(ctx.SourceType)
        if len(params) > 0 {
            res = strings.ReplaceAll(res, "{length}", params[0])
        }
    }
    return res
}

func toInt(v any) (int, bool) {
    switch val := v.(type) {
    case int: return val, true
    case float64: return int(val), true
    default: return 0, false
    }
}

// duplicated helper (move to utils if possible, but keeping here for self-contained loader)
func parseTypeWithParams(fullType string) (string, []string) {
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
