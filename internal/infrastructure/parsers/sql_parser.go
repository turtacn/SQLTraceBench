package parsers

import (
	"regexp"
	"strings"
)

type SQLParser interface {
	Parse(query string) (*ParseResult, error)
	ListTables(ast string) []string
}

type ParseResult struct {
	Tables        []string
	QueryTemplate string
	Params        []string
}

type BaseParser struct{}

func (bp *BaseParser) ListTables(ast string) (tables []string) {
	// very basic regex table extraction
	re := regexp.MustCompile(`(?i)(?:FROM|JOIN)\s+([^\s(,\)]+)`)
	matches := re.FindAllStringSubmatch(ast, -1)
	for _, m := range matches {
		tbl := strings.Trim(m[1], "`\"'")
		tables = append(tables, tbl)
	}
	return unique(tables)
}

func unique(s []string) []string {
	m := map[string]bool{}
	for _, v := range s {
		m[v] = true
	}
	var res []string
	for k := range m {
		res = append(res, k)
	}
	return res
}
//Personal.AI order the ending
