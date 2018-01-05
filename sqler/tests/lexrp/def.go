package lexrp

// Basic lexer, parser will do the rest
const ( // Lexical thing
	ItemError ItemType = iota
	ItemEOF
	ItemIdent // Ident?
	ItemComma
	ItemDoubleQuote
	ItemSingleQuote
	ItemLParen
	ItemRParen
	ItemUnknown
)

var (
	TokenStr = map[ItemType]string{
		ItemError:       "ERR",
		ItemIdent:       "Ident",
		ItemSingleQuote: "'string'",
		ItemDoubleQuote: "\"string\"",
		ItemLParen:      "open Parentesis",
		ItemRParen:      "close Parentesis",
		ItemUnknown:     "unknown",
		//EQUAL:    "Equal Operator",
		//ASTERISK: "Asterisk",
	}
)
