package parse

import (
	"fmt"
)

// Iter might be the walker
// State ful iterator that can be passed between functions
type Iter struct {
	tokens []*token
	pos    int
}

// Peek returns token at the current offset
func (it *Iter) Peek() *token {
	if it.pos < 0 || it.pos >= len(it.tokens) {
		return nil
	}
	return it.tokens[it.pos]
}

// Next return token and increases offset, moves next
func (it *Iter) Next() *token {
	it.pos++
	if it.pos >= len(it.tokens) {
		return nil
	}
	return it.tokens[it.pos]
}

// Clone returns a copy of Iter in current state
func (it *Iter) Clone() *Iter {
	return &Iter{it.tokens, it.pos}
}

// Reset or Rewind offset = 0,
func (it *Iter) Reset() *Iter {
	it.pos = 0
	return it // self
}

// String just a stringer to print the iterator
func (it *Iter) String() string {
	if it.pos >= len(it.tokens)-1 {
		return fmt.Sprintf("Nothign to see here, token is nil")
	}
	return fmt.Sprintf("Pos: %d, Rem: %d", it.pos, len(it.tokens[it.pos:]))
}
