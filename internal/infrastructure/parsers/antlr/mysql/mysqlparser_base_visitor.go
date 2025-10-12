// Code generated from MySqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // MySqlParser

import "github.com/antlr4-go/antlr/v4"

type BaseMySqlParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseMySqlParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMySqlParserVisitor) VisitSelect_statement(ctx *Select_statementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMySqlParserVisitor) VisitJoin_clause(ctx *Join_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMySqlParserVisitor) VisitTable_reference(ctx *Table_referenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMySqlParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMySqlParserVisitor) VisitAtom(ctx *AtomContext) interface{} {
	return v.VisitChildren(ctx)
}
