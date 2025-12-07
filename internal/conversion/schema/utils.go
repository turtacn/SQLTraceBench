package schema

import (
    "regexp"
    "strings"
)

// SplitWithBalance splits a string by a separator, respecting parenthesis balance.
func SplitWithBalance(s string, sep rune) []string {
	var parts []string
	var current strings.Builder
	balance := 0
	quote := rune(0)

	for _, r := range s {
		if quote != 0 {
			if r == quote {
				quote = 0
			}
			current.WriteRune(r)
			continue
		}

		if r == '\'' || r == '"' || r == '`' {
			quote = r
			current.WriteRune(r)
			continue
		}

		if r == '(' {
			balance++
		} else if r == ')' {
			balance--
		}

		if r == sep && balance == 0 {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

// ParseTypeWithParams parses a type string into base type and parameters.
// e.g. "DECIMAL(10, 2)" -> "DECIMAL", ["10", "2"]
func ParseTypeWithParams(fullType string) (string, []string) {
	re := regexp.MustCompile(`^([a-zA-Z0-9_ ]+)(?:\(([^)]+)\))?.*$`)
	matches := re.FindStringSubmatch(fullType)
	if len(matches) < 2 {
		return fullType, nil
	}
	baseType := strings.TrimSpace(matches[1])
	var params []string
	if len(matches) > 2 && matches[2] != "" {
		rawParams := SplitWithBalance(matches[2], ',') // Use SplitWithBalance for params too
		for _, p := range rawParams {
			params = append(params, strings.TrimSpace(p))
		}
	}
	return baseType, params
}
