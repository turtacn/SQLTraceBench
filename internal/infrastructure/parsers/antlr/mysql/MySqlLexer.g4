// MySqlLexer.g4
lexer grammar MySqlLexer;

// Keywords
SELECT: S E L E C T;
FROM: F R O M;
WHERE: W H E R E;
AND: A N D;
OR: O R;
NOT: N O T;
NULL_LITERAL: N U L L;
JOIN: J O I N;
ON: O N;

// Identifiers
ID: [a-zA-Z_] [a-zA-Z_0-9]*;
QUOTED_ID: '`' ( '``' | ~'`' )* '`';

// Literals
INT: [0-9]+;
STRING: '\'' ('\'\'' | ~'\'')* '\'';

// Operators
EQ: '=';
NEQ: '!=' | '<>';
GT: '>';
LT: '<';
GTE: '>=';
LTE: '<=';
STAR: '*';

// Punctuation
LPAREN: '(';
RPAREN: ')';
COMMA: ',';
DOT: '.';
SEMICOLON: ';';

// Whitespace
WS: [ \t\r\n]+ -> skip;

// Case-insensitive fragments
fragment A: [aA];
fragment B: [bB];
fragment C: [cC];
fragment D: [dD];
fragment E: [eE];
fragment F: [fF];
fragment G: [gG];
fragment H: [hH];
fragment I: [iI];
fragment J: [jJ];
fragment K: [kK];
fragment L: [lL];
fragment M: [mM];
fragment N: [nN];
fragment O: [oO];
fragment P: [pP];
fragment Q: [qQ];
fragment R: [rR];
fragment S: [sS];
fragment T: [tT];
fragment U: [uU];
fragment V: [vV];
fragment W: [wW];
fragment X: [xX];
fragment Y: [yY];
fragment Z: [zZ];