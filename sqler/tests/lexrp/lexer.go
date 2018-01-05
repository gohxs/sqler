package lexrp

import (
	"strings"
	"unicode/utf8"
)

const (
	eof = -1
)

type Lexer struct {
	input string
	start int
	pos   int

	width int // current rune width

	state stateFn
	items chan Item // we will send this there

}

// Start lexer
func Lex(input string) (*Lexer, chan Item) {
	l := &Lexer{
		input: input,
		state: lexBase,
		items: make(chan Item, 2),
	}
	go l.run()
	return l, l.items
}

func (l *Lexer) emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos], l.start}
	l.start = l.pos // next
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore a portion of the text
func (l *Lexer) ignore() {
	l.start = l.pos
}

// Go back the last readed rune by its width
// ideally the width should be recalculated to the previous utf size
func (l *Lexer) backup() {
	l.pos -= l.width
}

// Peek current rune from the string without moving cursor
func (l *Lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

// Utils for lexer?!
//

// Accept any character on the string
// not consuming it and returning false
func (l *Lexer) acceptAny(any string) bool {
	if strings.IndexRune(any, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// move the cursor forward while the char matches
// any on the string
func (l *Lexer) acceptAnyRun(any string) {
	for strings.IndexRune(any, l.next()) >= 0 {
	}
	l.backup()
}

// Run will move through states and close channel after
func (l *Lexer) run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// will fetch next item from channel or will process
// other character state
//
func (l *Lexer) NextItem() Item {
	item := <-l.items
	return item

	/*for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}*/
}

///////////////////////
/// State handlers

// STATES I GUESS
func lexBase(l *Lexer) stateFn {
	for {
		ch := l.next()
		switch ch {
		case ' ', '\t', '\r', '\f':
			l.ignore()
			continue // next
		case '(':
			l.emit(ItemLParen)
			return lexBase
		case ')':
			l.emit(ItemRParen)
			return lexBase
		case eof:
			return nil
		case '"':
			return lexDoubleQuotes
		case '\'':
			return lexSingleQuotes
		case '=', '*':
			l.ignore()
			return lexBase
		default:
			if isIdentStart(ch) {
				l.backup()
				return lexIdent
			}
		}
		l.emit(ItemUnknown)
		return lexBase
	}
}

// Atom
func lexIdent(l *Lexer) stateFn {
	for isIdentMiddle(l.next()) {
	}
	l.backup() // Back the last rune
	l.emit(ItemIdent)

	// Whats next??!
	return lexBase
}
func lexSingleQuotes(l *Lexer) stateFn {
	vi := strings.Index(l.input[l.pos:], "'")
	if vi == -1 {
		l.emit(ItemError)
		return nil
	}
	l.pos += vi + 1 // include '
	l.emit(ItemSingleQuote)
	return lexBase

}
func lexDoubleQuotes(l *Lexer) stateFn {
	vi := strings.Index(l.input[l.pos:], "\"")
	if vi == -1 {
		l.emit(ItemError)
		return nil
	}
	l.pos += vi + 1 // Include "
	l.emit(ItemDoubleQuote)
	return lexBase // back
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\f'
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

type stateFn func(*Lexer) stateFn

// State functions
