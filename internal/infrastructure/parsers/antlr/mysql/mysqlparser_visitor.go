// Code generated from MySqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // MySqlParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by MySqlParser.
type MySqlParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by MySqlParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by MySqlParser#select_statement.
	VisitSelect_statement(ctx *Select_statementContext) interface{}

	// Visit a parse tree produced by MySqlParser#join_clause.
	VisitJoin_clause(ctx *Join_clauseContext) interface{}

	// Visit a parse tree produced by MySqlParser#table_reference.
	VisitTable_reference(ctx *Table_referenceContext) interface{}

	// Visit a parse tree produced by MySqlParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by MySqlParser#atom.
	VisitAtom(ctx *AtomContext) interface{}
}
