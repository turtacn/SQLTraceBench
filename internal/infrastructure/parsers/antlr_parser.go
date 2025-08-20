package parsers

type ANTLRSQLParser struct {
	*BaseParser
}

func NewANTLRSQLParser() *ANTLRSQLParser {
	return &ANTLRSQLParser{BaseParser: &BaseParser{}}
}

//Personal.AI order the ending
