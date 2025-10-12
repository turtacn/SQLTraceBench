package parsers

import (
	"github.com/antlr4-go/antlr/v4"
	mysql "github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers/antlr/mysql"
)

// AntlrParser is a SQL parser that uses ANTLR to build a parse tree.
type AntlrParser struct{}

// NewAntlrParser creates a new AntlrParser.
func NewAntlrParser() *AntlrParser {
	return &AntlrParser{}
}

// ListTables extracts table names from a SQL query using the ANTLR parser.
func (p *AntlrParser) ListTables(sql string) ([]string, error) {
	// Create the ANTLR input stream.
	is := antlr.NewInputStream(sql)

	// Create the lexer.
	lexer := mysql.NewMySqlLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the parser.
	parser := mysql.NewMySqlParser(stream)
	parser.BuildParseTrees = true

	// Parse the query.
	tree := parser.Query()

	// Create the listener and walk the parse tree.
	listener := NewTableListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.TableNames, nil
}