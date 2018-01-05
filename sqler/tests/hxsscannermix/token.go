package fatihbased

const ( // Lexical thing
	TokenError TokenType = iota
	TokenEOF
	TokenIdent // Ident? Atom?
	TokenNumber
	TokenFloat
	TokenComma
	TokenDoubleQuote
	TokenSingleQuote
	TokenLParen
	TokenRParen
	TokenSymbol
	TokenUnknown
)

var (
	TokenStr = map[TokenType]string{
		TokenError:       "TokenError",
		TokenIdent:       "TokenIdent",
		TokenNumber:      "TokenNumber",
		TokenFloat:       "TokenFloat",
		TokenComma:       "TokenComma",
		TokenSingleQuote: "TokenSingleQuote",
		TokenDoubleQuote: "tokenDoubleQuote",
		TokenLParen:      "TokenLParen",
		TokenRParen:      "TokenRParen",
		TokenSymbol:      "TokenSymbol",
		TokenUnknown:     "TokenUnknown",
		//EQUAL:    "Equal Operator",
		//ASTERISK: "Asterisk",
	}
)

type TokenType int

type Token struct {
	Typ TokenType
	Val string
	Pos Pos // Extra position in the input string
}
