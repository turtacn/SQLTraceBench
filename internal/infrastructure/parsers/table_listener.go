package parsers

import mysql "github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers/antlr/mysql"

// TableListener is an ANTLR listener that extracts table names from a SQL query.
type TableListener struct {
	*mysql.BaseMySqlParserListener
	TableNames []string
}

// NewTableListener creates a new TableListener.
func NewTableListener() *TableListener {
	return &TableListener{
		BaseMySqlParserListener: &mysql.BaseMySqlParserListener{},
		TableNames:              make([]string, 0),
	}
}

// EnterTable_reference is called when the listener enters a `table_reference` node in the parse tree.
// It extracts the table name from the node and adds it to the list of table names.
func (l *TableListener) EnterTable_reference(ctx *mysql.Table_referenceContext) {
	l.TableNames = append(l.TableNames, ctx.GetText())
}