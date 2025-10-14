package parser

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/stretchr/testify/assert"
)

type testErrorListener struct {
	*antlr.DefaultErrorListener
	t      *testing.T
	errors int
}

func (l *testErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	l.errors++
}

func (l *testErrorListener) GetSyntaxErrors() int {
	return l.errors
}

func parseQuery(t *testing.T, sql string) (antlr.Tree, *testErrorListener) {
	t.Helper()

	// Create the lexer and parser
	input := antlr.NewInputStream(sql)
	lexer := NewMySqlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewMySqlParser(stream)

	// Add an error listener to capture syntax errors
	errorListener := &testErrorListener{t: t}
	p.AddErrorListener(errorListener)

	// Parse the query
	tree := p.Query()

	return tree, errorListener
}

func TestMySqlParser_SimpleSelect(t *testing.T) {
	// Parse a simple SELECT query
	tree, errorListener := parseQuery(t, "SELECT * FROM users;")

	// Assert that there are no syntax errors
	assert.Equal(t, 0, errorListener.GetSyntaxErrors(), "should not have any syntax errors")

	// Assert that the parse tree is not nil
	assert.NotNil(t, tree)
}

func TestMySqlParser_SelectWithWhere(t *testing.T) {
	// Parse a SELECT query with a WHERE clause
	tree, errorListener := parseQuery(t, "SELECT id, name FROM users WHERE id = 1;")

	// Assert that there are no syntax errors
	assert.Equal(t, 0, errorListener.GetSyntaxErrors(), "should not have any syntax errors")

	// Assert that the parse tree is not nil
	assert.NotNil(t, tree)
}

func TestMySqlParser_SelectWithJoin(t *testing.T) {
	// Parse a SELECT query with a JOIN clause
	tree, errorListener := parseQuery(t, "SELECT u.id, o.id FROM users u JOIN orders o ON u.id = o.user_id;")

	// Assert that there are no syntax errors
	assert.Equal(t, 0, errorListener.GetSyntaxErrors(), "should not have any syntax errors")

	// Assert that the parse tree is not nil
	assert.NotNil(t, tree)
}

func TestMySqlParser_InvalidQuery(t *testing.T) {
	// Parse an invalid query
	_, errorListener := parseQuery(t, "SELECT * FROM;")

	// Assert that there is at least one syntax error
	assert.Greater(t, errorListener.GetSyntaxErrors(), 0, "should have at least one syntax error")
}