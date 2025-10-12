// Code generated from MySqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // MySqlParser

import "github.com/antlr4-go/antlr/v4"

// BaseMySqlParserListener is a complete listener for a parse tree produced by MySqlParser.
type BaseMySqlParserListener struct{}

var _ MySqlParserListener = &BaseMySqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseMySqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseMySqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseMySqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseMySqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseMySqlParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseMySqlParserListener) ExitQuery(ctx *QueryContext) {}

// EnterSelect_statement is called when production select_statement is entered.
func (s *BaseMySqlParserListener) EnterSelect_statement(ctx *Select_statementContext) {}

// ExitSelect_statement is called when production select_statement is exited.
func (s *BaseMySqlParserListener) ExitSelect_statement(ctx *Select_statementContext) {}

// EnterJoin_clause is called when production join_clause is entered.
func (s *BaseMySqlParserListener) EnterJoin_clause(ctx *Join_clauseContext) {}

// ExitJoin_clause is called when production join_clause is exited.
func (s *BaseMySqlParserListener) ExitJoin_clause(ctx *Join_clauseContext) {}

// EnterTable_reference is called when production table_reference is entered.
func (s *BaseMySqlParserListener) EnterTable_reference(ctx *Table_referenceContext) {}

// ExitTable_reference is called when production table_reference is exited.
func (s *BaseMySqlParserListener) ExitTable_reference(ctx *Table_referenceContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseMySqlParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseMySqlParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterAtom is called when production atom is entered.
func (s *BaseMySqlParserListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BaseMySqlParserListener) ExitAtom(ctx *AtomContext) {}
