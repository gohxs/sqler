package lexer

var eof = rune(-1)

//type Rules struct {
//	Rule("select", OneOrMore(Or(fnIdent,fnString), ","), "from", fnDynTable, "where", fnIdent)
//}

const (
	EOF = iota
	UNKNOWN
	WHITESPACE
	IDENT
	STRING
	SINGLEQUOTESTR
	COMMA
	DOT
	// SYMBOLS

	EQUAL
	ASTERISK
)

var (
	TokenStr = map[int]string{
		UNKNOWN:        "Unknown",
		WHITESPACE:     "Whitespace",
		IDENT:          "Ident",
		STRING:         "String",
		SINGLEQUOTESTR: "Single String",
		COMMA:          "Comma",
		DOT:            "Dot",
		//
		EQUAL:    "Equal Operator",
		ASTERISK: "Asterisk",
	}
)

type synToken struct {
	id  int
	pos int
	str string
	// Content tokens if it is a block
}

type Scanner struct {
	in       string // Rune string
	pos      int
	curToken synToken
	// rules and rule matcher
	///
	//current token/suggestion
}

func NewScanner(in string) *Scanner {
	return &Scanner{in, 0, synToken{}}
}

func (s *Scanner) peek() rune { // Rune like?
	if s.pos >= len(s.in) {
		return eof
	}
	return rune(s.in[s.pos])
}
func (s *Scanner) next() rune {
	ch := s.peek()
	if ch != eof {
		s.pos++
	}
	return ch
}

// Scan Fetch a token or suggestion?
func (s *Scanner) scan(lval *synToken) {
	s.skipWhiteSpace()

	lval.id = EOF
	lval.pos = s.pos // Starting position
	lval.str = "EOF"

	ch := s.next()
	if ch == eof {
		return

	}
	lval.id = UNKNOWN

	lval.pos = s.pos - 1
	lval.str = s.in[lval.pos:s.pos]

	switch ch {
	case '"':
		s.scanString(lval, ch) // Scan string to ch
	case '\'':
		s.scanString(lval, ch) // Scan string to ch
		lval.id = SINGLEQUOTESTR
	case '=':
		lval.id = EQUAL
	case '*':
		lval.id = ASTERISK
	default:
		if isIdentStart(ch) {
			s.scanIdent(lval)
		}
	}

	return
}
func (s *Scanner) skipWhiteSpace() {
	for {
		ch := s.peek()
		if ch == ' ' || ch == '\t' || ch == '\f' || ch == '\f' {
			s.pos++
			continue
		}
		break
	}
}
func (s *Scanner) scanIdent(lval *synToken) {
	start := s.pos - 1
	for ; isIdentMiddle(s.peek()); s.pos++ {
	}
	lval.str = s.in[start:s.pos]
	lval.id = IDENT
}

// Read until
func (s *Scanner) scanString(lval *synToken, end rune) {
	start := s.pos
	for {
		ch := s.next()
		if ch == end && s.in[s.pos-1] != '\\' || ch == eof {
			break
		}
	}
	lval.str = s.in[start : s.pos-1]
	// Define keyworder
	lval.id = STRING
}

func isIdentStart(ch rune) bool {
	return (ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' ||
		ch == '_')
}
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
func isIdentMiddle(ch rune) bool {
	return isIdentStart(ch) || isDigit(ch)
}
