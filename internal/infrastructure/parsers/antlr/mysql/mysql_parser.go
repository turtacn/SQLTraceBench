// Code generated from MySqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // MySqlParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type MySqlParser struct {
	*antlr.BaseParser
}

var MySqlParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func mysqlparserParserInit() {
	staticData := &MySqlParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "'='", "", "'>'",
		"'<'", "'>='", "'<='", "'*'", "'('", "')'", "','", "'.'", "';'",
	}
	staticData.SymbolicNames = []string{
		"", "SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "NULL_LITERAL", "JOIN",
		"ON", "ID", "QUOTED_ID", "INT", "STRING", "EQ", "NEQ", "GT", "LT", "GTE",
		"LTE", "STAR", "LPAREN", "RPAREN", "COMMA", "DOT", "SEMICOLON", "WS",
	}
	staticData.RuleNames = []string{
		"query", "select_statement", "join_clause", "table_reference", "expression",
		"atom",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 26, 81, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 21,
		8, 1, 10, 1, 12, 1, 24, 9, 1, 3, 1, 26, 8, 1, 1, 1, 1, 1, 1, 1, 5, 1, 31,
		8, 1, 10, 1, 12, 1, 34, 9, 1, 1, 1, 1, 1, 3, 1, 38, 8, 1, 1, 2, 1, 2, 1,
		2, 1, 2, 3, 2, 44, 8, 2, 1, 3, 1, 3, 3, 3, 48, 8, 3, 1, 4, 1, 4, 1, 4,
		1, 4, 3, 4, 54, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 62, 8,
		4, 10, 4, 12, 4, 65, 9, 4, 1, 5, 1, 5, 1, 5, 3, 5, 70, 8, 5, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 79, 8, 5, 1, 5, 0, 1, 8, 6, 0, 2, 4,
		6, 8, 10, 0, 3, 1, 0, 10, 11, 1, 0, 4, 5, 1, 0, 14, 19, 88, 0, 12, 1, 0,
		0, 0, 2, 15, 1, 0, 0, 0, 4, 39, 1, 0, 0, 0, 6, 45, 1, 0, 0, 0, 8, 53, 1,
		0, 0, 0, 10, 78, 1, 0, 0, 0, 12, 13, 3, 2, 1, 0, 13, 14, 5, 25, 0, 0, 14,
		1, 1, 0, 0, 0, 15, 25, 5, 1, 0, 0, 16, 26, 5, 20, 0, 0, 17, 22, 3, 8, 4,
		0, 18, 19, 5, 23, 0, 0, 19, 21, 3, 8, 4, 0, 20, 18, 1, 0, 0, 0, 21, 24,
		1, 0, 0, 0, 22, 20, 1, 0, 0, 0, 22, 23, 1, 0, 0, 0, 23, 26, 1, 0, 0, 0,
		24, 22, 1, 0, 0, 0, 25, 16, 1, 0, 0, 0, 25, 17, 1, 0, 0, 0, 26, 27, 1,
		0, 0, 0, 27, 28, 5, 2, 0, 0, 28, 32, 3, 6, 3, 0, 29, 31, 3, 4, 2, 0, 30,
		29, 1, 0, 0, 0, 31, 34, 1, 0, 0, 0, 32, 30, 1, 0, 0, 0, 32, 33, 1, 0, 0,
		0, 33, 37, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 35, 36, 5, 3, 0, 0, 36, 38,
		3, 8, 4, 0, 37, 35, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0, 38, 3, 1, 0, 0, 0,
		39, 40, 5, 8, 0, 0, 40, 43, 3, 6, 3, 0, 41, 42, 5, 9, 0, 0, 42, 44, 3,
		8, 4, 0, 43, 41, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44, 5, 1, 0, 0, 0, 45,
		47, 7, 0, 0, 0, 46, 48, 7, 0, 0, 0, 47, 46, 1, 0, 0, 0, 47, 48, 1, 0, 0,
		0, 48, 7, 1, 0, 0, 0, 49, 50, 6, 4, -1, 0, 50, 54, 3, 10, 5, 0, 51, 52,
		5, 6, 0, 0, 52, 54, 3, 8, 4, 2, 53, 49, 1, 0, 0, 0, 53, 51, 1, 0, 0, 0,
		54, 63, 1, 0, 0, 0, 55, 56, 10, 3, 0, 0, 56, 57, 7, 1, 0, 0, 57, 62, 3,
		8, 4, 4, 58, 59, 10, 1, 0, 0, 59, 60, 7, 2, 0, 0, 60, 62, 3, 8, 4, 2, 61,
		55, 1, 0, 0, 0, 61, 58, 1, 0, 0, 0, 62, 65, 1, 0, 0, 0, 63, 61, 1, 0, 0,
		0, 63, 64, 1, 0, 0, 0, 64, 9, 1, 0, 0, 0, 65, 63, 1, 0, 0, 0, 66, 69, 7,
		0, 0, 0, 67, 68, 5, 24, 0, 0, 68, 70, 7, 0, 0, 0, 69, 67, 1, 0, 0, 0, 69,
		70, 1, 0, 0, 0, 70, 79, 1, 0, 0, 0, 71, 79, 5, 12, 0, 0, 72, 79, 5, 13,
		0, 0, 73, 79, 5, 7, 0, 0, 74, 75, 5, 21, 0, 0, 75, 76, 3, 8, 4, 0, 76,
		77, 5, 22, 0, 0, 77, 79, 1, 0, 0, 0, 78, 66, 1, 0, 0, 0, 78, 71, 1, 0,
		0, 0, 78, 72, 1, 0, 0, 0, 78, 73, 1, 0, 0, 0, 78, 74, 1, 0, 0, 0, 79, 11,
		1, 0, 0, 0, 11, 22, 25, 32, 37, 43, 47, 53, 61, 63, 69, 78,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// MySqlParserInit initializes any static state used to implement MySqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewMySqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func MySqlParserInit() {
	staticData := &MySqlParserParserStaticData
	staticData.once.Do(mysqlparserParserInit)
}

// NewMySqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewMySqlParser(input antlr.TokenStream) *MySqlParser {
	MySqlParserInit()
	this := new(MySqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &MySqlParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "MySqlParser.g4"

	return this
}

// MySqlParser tokens.
const (
	MySqlParserEOF          = antlr.TokenEOF
	MySqlParserSELECT       = 1
	MySqlParserFROM         = 2
	MySqlParserWHERE        = 3
	MySqlParserAND          = 4
	MySqlParserOR           = 5
	MySqlParserNOT          = 6
	MySqlParserNULL_LITERAL = 7
	MySqlParserJOIN         = 8
	MySqlParserON           = 9
	MySqlParserID           = 10
	MySqlParserQUOTED_ID    = 11
	MySqlParserINT          = 12
	MySqlParserSTRING       = 13
	MySqlParserEQ           = 14
	MySqlParserNEQ          = 15
	MySqlParserGT           = 16
	MySqlParserLT           = 17
	MySqlParserGTE          = 18
	MySqlParserLTE          = 19
	MySqlParserSTAR         = 20
	MySqlParserLPAREN       = 21
	MySqlParserRPAREN       = 22
	MySqlParserCOMMA        = 23
	MySqlParserDOT          = 24
	MySqlParserSEMICOLON    = 25
	MySqlParserWS           = 26
)

// MySqlParser rules.
const (
	MySqlParserRULE_query            = 0
	MySqlParserRULE_select_statement = 1
	MySqlParserRULE_join_clause      = 2
	MySqlParserRULE_table_reference  = 3
	MySqlParserRULE_expression       = 4
	MySqlParserRULE_atom             = 5
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Select_statement() ISelect_statementContext
	SEMICOLON() antlr.TerminalNode

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) Select_statement() ISelect_statementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelect_statementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelect_statementContext)
}

func (s *QueryContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(MySqlParserSEMICOLON, 0)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (s *QueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, MySqlParserRULE_query)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(12)
		p.Select_statement()
	}
	{
		p.SetState(13)
		p.Match(MySqlParserSEMICOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelect_statementContext is an interface to support dynamic dispatch.
type ISelect_statementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELECT() antlr.TerminalNode
	FROM() antlr.TerminalNode
	Table_reference() ITable_referenceContext
	STAR() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllJoin_clause() []IJoin_clauseContext
	Join_clause(i int) IJoin_clauseContext
	WHERE() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSelect_statementContext differentiates from other interfaces.
	IsSelect_statementContext()
}

type Select_statementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelect_statementContext() *Select_statementContext {
	var p = new(Select_statementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_select_statement
	return p
}

func InitEmptySelect_statementContext(p *Select_statementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_select_statement
}

func (*Select_statementContext) IsSelect_statementContext() {}

func NewSelect_statementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Select_statementContext {
	var p = new(Select_statementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_select_statement

	return p
}

func (s *Select_statementContext) GetParser() antlr.Parser { return s.parser }

func (s *Select_statementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(MySqlParserSELECT, 0)
}

func (s *Select_statementContext) FROM() antlr.TerminalNode {
	return s.GetToken(MySqlParserFROM, 0)
}

func (s *Select_statementContext) Table_reference() ITable_referenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITable_referenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITable_referenceContext)
}

func (s *Select_statementContext) STAR() antlr.TerminalNode {
	return s.GetToken(MySqlParserSTAR, 0)
}

func (s *Select_statementContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Select_statementContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Select_statementContext) AllJoin_clause() []IJoin_clauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IJoin_clauseContext); ok {
			len++
		}
	}

	tst := make([]IJoin_clauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IJoin_clauseContext); ok {
			tst[i] = t.(IJoin_clauseContext)
			i++
		}
	}

	return tst
}

func (s *Select_statementContext) Join_clause(i int) IJoin_clauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoin_clauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoin_clauseContext)
}

func (s *Select_statementContext) WHERE() antlr.TerminalNode {
	return s.GetToken(MySqlParserWHERE, 0)
}

func (s *Select_statementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(MySqlParserCOMMA)
}

func (s *Select_statementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(MySqlParserCOMMA, i)
}

func (s *Select_statementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Select_statementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Select_statementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterSelect_statement(s)
	}
}

func (s *Select_statementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitSelect_statement(s)
	}
}

func (s *Select_statementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitSelect_statement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Select_statement() (localctx ISelect_statementContext) {
	localctx = NewSelect_statementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, MySqlParserRULE_select_statement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(15)
		p.Match(MySqlParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(25)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case MySqlParserSTAR:
		{
			p.SetState(16)
			p.Match(MySqlParserSTAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case MySqlParserNOT, MySqlParserNULL_LITERAL, MySqlParserID, MySqlParserQUOTED_ID, MySqlParserINT, MySqlParserSTRING, MySqlParserLPAREN:
		{
			p.SetState(17)
			p.expression(0)
		}
		p.SetState(22)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == MySqlParserCOMMA {
			{
				p.SetState(18)
				p.Match(MySqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(19)
				p.expression(0)
			}

			p.SetState(24)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	{
		p.SetState(27)
		p.Match(MySqlParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(28)
		p.Table_reference()
	}
	p.SetState(32)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == MySqlParserJOIN {
		{
			p.SetState(29)
			p.Join_clause()
		}

		p.SetState(34)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == MySqlParserWHERE {
		{
			p.SetState(35)
			p.Match(MySqlParserWHERE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(36)
			p.expression(0)
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoin_clauseContext is an interface to support dynamic dispatch.
type IJoin_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	JOIN() antlr.TerminalNode
	Table_reference() ITable_referenceContext
	ON() antlr.TerminalNode
	Expression() IExpressionContext

	// IsJoin_clauseContext differentiates from other interfaces.
	IsJoin_clauseContext()
}

type Join_clauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoin_clauseContext() *Join_clauseContext {
	var p = new(Join_clauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_join_clause
	return p
}

func InitEmptyJoin_clauseContext(p *Join_clauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_join_clause
}

func (*Join_clauseContext) IsJoin_clauseContext() {}

func NewJoin_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Join_clauseContext {
	var p = new(Join_clauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_join_clause

	return p
}

func (s *Join_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Join_clauseContext) JOIN() antlr.TerminalNode {
	return s.GetToken(MySqlParserJOIN, 0)
}

func (s *Join_clauseContext) Table_reference() ITable_referenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITable_referenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITable_referenceContext)
}

func (s *Join_clauseContext) ON() antlr.TerminalNode {
	return s.GetToken(MySqlParserON, 0)
}

func (s *Join_clauseContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Join_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Join_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Join_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterJoin_clause(s)
	}
}

func (s *Join_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitJoin_clause(s)
	}
}

func (s *Join_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitJoin_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Join_clause() (localctx IJoin_clauseContext) {
	localctx = NewJoin_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, MySqlParserRULE_join_clause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(39)
		p.Match(MySqlParserJOIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(40)
		p.Table_reference()
	}
	p.SetState(43)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == MySqlParserON {
		{
			p.SetState(41)
			p.Match(MySqlParserON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(42)
			p.expression(0)
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITable_referenceContext is an interface to support dynamic dispatch.
type ITable_referenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllQUOTED_ID() []antlr.TerminalNode
	QUOTED_ID(i int) antlr.TerminalNode

	// IsTable_referenceContext differentiates from other interfaces.
	IsTable_referenceContext()
}

type Table_referenceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTable_referenceContext() *Table_referenceContext {
	var p = new(Table_referenceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_table_reference
	return p
}

func InitEmptyTable_referenceContext(p *Table_referenceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_table_reference
}

func (*Table_referenceContext) IsTable_referenceContext() {}

func NewTable_referenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Table_referenceContext {
	var p = new(Table_referenceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_table_reference

	return p
}

func (s *Table_referenceContext) GetParser() antlr.Parser { return s.parser }

func (s *Table_referenceContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(MySqlParserID)
}

func (s *Table_referenceContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(MySqlParserID, i)
}

func (s *Table_referenceContext) AllQUOTED_ID() []antlr.TerminalNode {
	return s.GetTokens(MySqlParserQUOTED_ID)
}

func (s *Table_referenceContext) QUOTED_ID(i int) antlr.TerminalNode {
	return s.GetToken(MySqlParserQUOTED_ID, i)
}

func (s *Table_referenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Table_referenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Table_referenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterTable_reference(s)
	}
}

func (s *Table_referenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitTable_reference(s)
	}
}

func (s *Table_referenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitTable_reference(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Table_reference() (localctx ITable_referenceContext) {
	localctx = NewTable_referenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, MySqlParserRULE_table_reference)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(45)
		_la = p.GetTokenStream().LA(1)

		if !(_la == MySqlParserID || _la == MySqlParserQUOTED_ID) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(47)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == MySqlParserID || _la == MySqlParserQUOTED_ID {
		{
			p.SetState(46)
			_la = p.GetTokenStream().LA(1)

			if !(_la == MySqlParserID || _la == MySqlParserQUOTED_ID) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Atom() IAtomContext
	NOT() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode
	EQ() antlr.TerminalNode
	NEQ() antlr.TerminalNode
	GT() antlr.TerminalNode
	LT() antlr.TerminalNode
	GTE() antlr.TerminalNode
	LTE() antlr.TerminalNode

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Atom() IAtomContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAtomContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *ExpressionContext) NOT() antlr.TerminalNode {
	return s.GetToken(MySqlParserNOT, 0)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) AND() antlr.TerminalNode {
	return s.GetToken(MySqlParserAND, 0)
}

func (s *ExpressionContext) OR() antlr.TerminalNode {
	return s.GetToken(MySqlParserOR, 0)
}

func (s *ExpressionContext) EQ() antlr.TerminalNode {
	return s.GetToken(MySqlParserEQ, 0)
}

func (s *ExpressionContext) NEQ() antlr.TerminalNode {
	return s.GetToken(MySqlParserNEQ, 0)
}

func (s *ExpressionContext) GT() antlr.TerminalNode {
	return s.GetToken(MySqlParserGT, 0)
}

func (s *ExpressionContext) LT() antlr.TerminalNode {
	return s.GetToken(MySqlParserLT, 0)
}

func (s *ExpressionContext) GTE() antlr.TerminalNode {
	return s.GetToken(MySqlParserGTE, 0)
}

func (s *ExpressionContext) LTE() antlr.TerminalNode {
	return s.GetToken(MySqlParserLTE, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *MySqlParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 8
	p.EnterRecursionRule(localctx, 8, MySqlParserRULE_expression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(53)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case MySqlParserNULL_LITERAL, MySqlParserID, MySqlParserQUOTED_ID, MySqlParserINT, MySqlParserSTRING, MySqlParserLPAREN:
		{
			p.SetState(50)
			p.Atom()
		}

	case MySqlParserNOT:
		{
			p.SetState(51)
			p.Match(MySqlParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(52)
			p.expression(2)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(61)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, MySqlParserRULE_expression)
				p.SetState(55)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(56)
					_la = p.GetTokenStream().LA(1)

					if !(_la == MySqlParserAND || _la == MySqlParserOR) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(57)
					p.expression(4)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, MySqlParserRULE_expression)
				p.SetState(58)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(59)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1032192) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(60)
					p.expression(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(65)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAtomContext is an interface to support dynamic dispatch.
type IAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllQUOTED_ID() []antlr.TerminalNode
	QUOTED_ID(i int) antlr.TerminalNode
	DOT() antlr.TerminalNode
	INT() antlr.TerminalNode
	STRING() antlr.TerminalNode
	NULL_LITERAL() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsAtomContext differentiates from other interfaces.
	IsAtomContext()
}

type AtomContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtomContext() *AtomContext {
	var p = new(AtomContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_atom
	return p
}

func InitEmptyAtomContext(p *AtomContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = MySqlParserRULE_atom
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = MySqlParserRULE_atom

	return p
}

func (s *AtomContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(MySqlParserID)
}

func (s *AtomContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(MySqlParserID, i)
}

func (s *AtomContext) AllQUOTED_ID() []antlr.TerminalNode {
	return s.GetTokens(MySqlParserQUOTED_ID)
}

func (s *AtomContext) QUOTED_ID(i int) antlr.TerminalNode {
	return s.GetToken(MySqlParserQUOTED_ID, i)
}

func (s *AtomContext) DOT() antlr.TerminalNode {
	return s.GetToken(MySqlParserDOT, 0)
}

func (s *AtomContext) INT() antlr.TerminalNode {
	return s.GetToken(MySqlParserINT, 0)
}

func (s *AtomContext) STRING() antlr.TerminalNode {
	return s.GetToken(MySqlParserSTRING, 0)
}

func (s *AtomContext) NULL_LITERAL() antlr.TerminalNode {
	return s.GetToken(MySqlParserNULL_LITERAL, 0)
}

func (s *AtomContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(MySqlParserLPAREN, 0)
}

func (s *AtomContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AtomContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(MySqlParserRPAREN, 0)
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.EnterAtom(s)
	}
}

func (s *AtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(MySqlParserListener); ok {
		listenerT.ExitAtom(s)
	}
}

func (s *AtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case MySqlParserVisitor:
		return t.VisitAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *MySqlParser) Atom() (localctx IAtomContext) {
	localctx = NewAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, MySqlParserRULE_atom)
	var _la int

	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case MySqlParserID, MySqlParserQUOTED_ID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(66)
			_la = p.GetTokenStream().LA(1)

			if !(_la == MySqlParserID || _la == MySqlParserQUOTED_ID) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(69)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(67)
				p.Match(MySqlParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(68)
				_la = p.GetTokenStream().LA(1)

				if !(_la == MySqlParserID || _la == MySqlParserQUOTED_ID) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case MySqlParserINT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(71)
			p.Match(MySqlParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case MySqlParserSTRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(72)
			p.Match(MySqlParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case MySqlParserNULL_LITERAL:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(73)
			p.Match(MySqlParserNULL_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case MySqlParserLPAREN:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(74)
			p.Match(MySqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(75)
			p.expression(0)
		}
		{
			p.SetState(76)
			p.Match(MySqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *MySqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 4:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *MySqlParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
