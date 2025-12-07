package schema

import "strings"

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
