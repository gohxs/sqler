package fatihbased

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Pos struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count)
}
type Error struct {
	Pos Pos
	Err error
}

// Scanner based on fatih HLC scanner
type Scanner struct {
	buf *bytes.Buffer
	src []byte

	// Positions
	srcPos  Pos
	prevPos Pos // also remembers last line

	// last character width
	// and last line width
	width       int
	lastLineLen int

	errors []Error

	//tokPos Pos // Why?

	curTok Token // The Current Token

}

const eof = rune(0) // maybe -1

// New creates a New Scanner
func New(src []byte) *Scanner {
	b := bytes.NewBuffer(src) // Swap this with a bufio reader perhaps
	s := &Scanner{
		buf: b,
		src: src, // Unfortunate
	}

	s.srcPos.Line = 1

	return s

}

// Based on fatih HCL scanner,
// Reducing code to common basis
func (s *Scanner) next() rune {
	prevPos := s.srcPos

	ch, size, err := s.buf.ReadRune()
	s.srcPos.Column++
	s.srcPos.Offset += size
	s.width = size

	// Err
	if err != nil {
		return eof
	}

	if ch == utf8.RuneError && size == 1 {
		s.errorf("illegal UTF-8 encoding")
		return ch
	}
	// use Remembered pos
	s.prevPos = prevPos

	if ch == '\n' {
		s.srcPos.Line++
		s.lastLineLen = s.srcPos.Column
		s.srcPos.Column = 0
	}
	return ch
}

// same as rob pike backup()
func (s *Scanner) unread() {
	if err := s.buf.UnreadRune(); err != nil {
		panic(err)
	}
	s.srcPos = s.prevPos
}

func (s *Scanner) peek() rune {
	peek, _, err := s.buf.ReadRune()
	if err != nil {
		return eof
	}
	s.buf.UnreadRune()

	return peek
}

// Mark a token from current position start/end
func (s *Scanner) token(t TokenType) {
	s.curTok.Typ = t
	s.curTok.Val = string(s.src[s.curTok.Pos.Offset:s.srcPos.Offset])
}

func (s *Scanner) errorf(format string, a ...interface{}) {
	s.errors = append(s.errors, Error{s.srcPos, fmt.Errorf(format, a...)})
}

// Token Return latest token
// Maybe returning pointer is not a good idea
func (s *Scanner) Token() *Token {
	return &s.curTok
}

// Scan token and save state to be fetched with
// (s *Scanner).Token()
func (s *Scanner) Scan() bool {
	ch := s.next()
	for isWhitespace(ch) {
		ch = s.next()
	}
	// Same as fatih
	s.curTok.Pos.Offset = s.srcPos.Offset - s.width // s.tokStart
	if s.srcPos.Column > 0 {
		// common case: last character was not a '\n'
		s.curTok.Pos.Line = s.srcPos.Line
		s.curTok.Pos.Column = s.srcPos.Column - 1 // ?
	} else {
		// last character was a '\n'
		// (we cannot be at the beginning of the source
		// since we have called next() at least once)
		s.curTok.Pos.Line = s.srcPos.Line - 1
		s.curTok.Pos.Column = s.lastLineLen
	}

	switch ch {
	case eof:
		s.token(TokenEOF)
		return false
	case '"':
		s.token(s.scanString('"'))
	case '\'':
		s.token(s.scanString('\''))
	case '(':
		s.token(TokenLParen)
	case ')':
		s.token(TokenRParen)
	case ',':
		s.token(TokenComma)
	case ';', '*': // etc
		s.token(TokenSymbol)
	default:
		switch {
		case isIdentStart(ch):
			s.token(s.scanIdentifier())
		case isDecimal(ch):
			s.token(s.scanNumber(ch)) // Should return Type here at least
		default:
			s.errorf("Unknown token")
			s.token(TokenUnknown)
		}
	}
	return true // Keep going
}

func (s *Scanner) scanIdentifier() TokenType {
	ch := s.next()
	for isLetter(ch) || isDigit(ch) {
		ch = s.next()
	}
	if ch != eof {
		s.unread()
	}

	return TokenIdent

}

// scanString scans a quoted string
func (s *Scanner) scanString(quoteMark rune) TokenType {
	for {
		// read character after quote
		ch := s.next()

		if ch == '\n' || ch < 0 || ch == eof {
			s.errorf("literal not terminated")
			return TokenError
		}

		if ch == quoteMark {
			return TokenDoubleQuote
		}

		if ch == '\\' {
			s.scanEscape()
		}
	}

	// unterminated too
	return TokenError
}

// scanEscape scans an escape sequence
// Why should we scan the escapes?
func (s *Scanner) scanEscape() rune {
	// http://en.cppreference.com/w/cpp/language/escape
	ch := s.next() // read character after '/'
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
		// nothing to do
	case '0', '1', '2', '3', '4', '5', '6', '7':
		// octal notation
		ch = s.scanDigits(ch, 8, 3)
	case 'x':
		// hexademical notation
		ch = s.scanDigits(s.next(), 16, 2)
	case 'u':
		// universal character name
		ch = s.scanDigits(s.next(), 16, 4)
	case 'U':
		// universal character name
		ch = s.scanDigits(s.next(), 16, 8)
	default:
		s.errorf("illegal char escape")
	}
	return ch
}
func (s *Scanner) scanDigits(ch rune, base, n int) rune {
	for n > 0 && digitVal(ch) < base {
		ch = s.next()
		n--
	}
	if n > 0 {
		s.errorf("illegal char escape")
	}

	// we scanned all digits, put the last non digit char back
	s.unread()
	return ch
}

// scanNumber scans a HCL number definition starting with the given rune
func (s *Scanner) scanNumber(ch rune) TokenType {
	if ch == '0' {
		// check for hexadecimal, octal or float
		ch = s.next()
		if ch == 'x' || ch == 'X' {
			// hexadecimal
			ch = s.next()
			found := false
			for isHexadecimal(ch) {
				ch = s.next()
				found = true
			}
			if !found {
				s.errorf("illegal hexadecimal number")
			}

			if ch != eof {
				s.unread()
			}

			return TokenNumber
		}

		// now it's either something like: 0421(octal) or 0.1231(float)
		illegalOctal := false
		for isDecimal(ch) {
			ch = s.next()
			if ch == '8' || ch == '9' {
				// this is just a possibility. For example 0159 is illegal, but
				// 0159.23 is valid. So we mark a possible illegal octal. If
				// the next character is not a period, we'll print the error.
				illegalOctal = true
			}
		}

		// literals of form 01e10 are treates as Numbers in HCL, which differs from Go.
		if ch == 'e' || ch == 'E' {
			ch = s.scanExponent(ch)
			return TokenNumber
		}

		if ch == '.' {
			ch = s.scanFraction(ch)

			if ch == 'e' || ch == 'E' {
				ch = s.next()
				ch = s.scanExponent(ch)
			}
			return TokenFloat
		}

		if illegalOctal {
			s.errorf("illegal octal number")
		}

		if ch != eof {
			s.unread()
		}
		return TokenNumber
	}

	s.scanMantissa(ch)
	ch = s.next() // seek forward
	// literals of form 1e10 are treates as Numbers in HCL, which differs from Go.
	if ch == 'e' || ch == 'E' {
		ch = s.scanExponent(ch)
		return TokenNumber
	}

	if ch == '.' {
		ch = s.scanFraction(ch)
		if ch == 'e' || ch == 'E' {
			ch = s.next()
			ch = s.scanExponent(ch)
		}
		return TokenFloat
	}

	if ch != eof {
		s.unread()
	}
	return TokenNumber
}

// Sub scans
// scanMantissa scans the mantissa begining from the rune. It returns the next
// non decimal rune. It's used to determine wheter it's a fraction or exponent.
func (s *Scanner) scanMantissa(ch rune) rune {
	scanned := false
	for isDecimal(ch) {
		ch = s.next()
		scanned = true
	}

	if scanned && ch != eof {
		s.unread()
	}
	return ch
}

// scanFraction scans the fraction after the '.' rune
func (s *Scanner) scanFraction(ch rune) rune {
	if ch == '.' {
		ch = s.peek() // we peek just to see if we can move forward
		ch = s.scanMantissa(ch)
	}
	return ch
}

// scanExponent scans the remaining parts of an exponent after the 'e' or 'E'
// rune.
func (s *Scanner) scanExponent(ch rune) rune {
	if ch == 'e' || ch == 'E' {
		ch = s.next()
		if ch == '-' || ch == '+' {
			ch = s.next()
		}
		ch = s.scanMantissa(ch)
	}
	return ch
}

// Identifiers

func isIdentStart(ch rune) bool {
	return isLetter(ch) || ch == '_'
}

// isHexadecimal returns true if the given rune is a letter
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
}

// isHexadecimal returns true if the given rune is a decimal digit
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
}

// isHexadecimal returns true if the given rune is a decimal number
func isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isHexadecimal returns true if the given rune is an hexadecimal number
func isHexadecimal(ch rune) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

// isWhitespace returns true if the rune is a space, tab, newline or carriage return
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}
