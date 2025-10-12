// MySqlParser.g4
parser grammar MySqlParser;

options { tokenVocab = MySqlLexer; }

// Rules
query: select_statement SEMICOLON;

select_statement:
    SELECT (STAR | expression (COMMA expression)*) FROM table_reference (join_clause)* (WHERE expression)?;

join_clause:
    JOIN table_reference (ON expression)?;

table_reference:
    (ID | QUOTED_ID) (ID | QUOTED_ID)?;

expression:
    atom
    | expression (AND | OR) expression
    | NOT expression
    | expression (EQ | NEQ | GT | LT | GTE | LTE) expression;

atom:
    (ID | QUOTED_ID) (DOT (ID | QUOTED_ID))?
    | INT
    | STRING
    | NULL_LITERAL
    | LPAREN expression RPAREN;