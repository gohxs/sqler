package scanner

const ( // Lexical thing
	ERROR TokenType = iota
	EOF
	IDENT // Ident? Atom?
	INT
	FLOAT
	COMMA
	DSTRING
	SSTRING
	LPAREN
	RPAREN
	SYMBOL
	UNKNOWN
)

var (
	tokenStr = map[TokenType]string{
		ERROR:   "ERROR",
		IDENT:   "IDENT",
		INT:     "INT",
		FLOAT:   "FLOAT",
		COMMA:   "COMMA",
		SSTRING: "SSTRING",
		DSTRING: "DSTRING",
		LPAREN:  "LPAREN",
		RPAREN:  "RPAREN",
		SYMBOL:  "SYMBOL",
		UNKNOWN: "UNKNOWN",
		//EQUAL:    "Equal Operator",
		//ASTERISK: "Asterisk",
	}
)

//TokenType the type of token
type TokenType int

func (t TokenType) String() string {
	return tokenStr[t]
}

/*type Token struct {
	Typ      TokenType
	Position Position // Extra position in the input string
	Val      func() string
}*/
