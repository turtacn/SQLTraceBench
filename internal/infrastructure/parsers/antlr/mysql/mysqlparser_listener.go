// Code generated from MySqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // MySqlParser

import "github.com/antlr4-go/antlr/v4"

// MySqlParserListener is a complete listener for a parse tree produced by MySqlParser.
type MySqlParserListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterSelect_statement is called when entering the select_statement production.
	EnterSelect_statement(c *Select_statementContext)

	// EnterJoin_clause is called when entering the join_clause production.
	EnterJoin_clause(c *Join_clauseContext)

	// EnterTable_reference is called when entering the table_reference production.
	EnterTable_reference(c *Table_referenceContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitSelect_statement is called when exiting the select_statement production.
	ExitSelect_statement(c *Select_statementContext)

	// ExitJoin_clause is called when exiting the join_clause production.
	ExitJoin_clause(c *Join_clauseContext)

	// ExitTable_reference is called when exiting the table_reference production.
	ExitTable_reference(c *Table_referenceContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)
}
