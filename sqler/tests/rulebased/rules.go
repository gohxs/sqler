//
//Lex completer
//
// The Idea:
//   the scanner go trough a set of rules,
//   it can only go to a next rule if fully matched the first
//   if it has complete match we move to next ruler
//
//   "S" - "SHOW","SELECT"
//   "SE" - "SELECT"
//   "SELECT *" - "SELECT", "*"
//   Engage sub scanners until EOF
//
//
package lexer

func Rule(r ...interface{}) {
	for _, k := range r {
		switch v := k.(type) {

		}
	}
}

var eof = rune(-1)

// Scanner the scanner struct
type Scanner struct {
	in  *string
	pos int
}

func (s *Scanner) peek() rune {
	if s.pos >= len(*s.in) {
		return eof
	}
	return rune((*s.in)[s.pos])
}

func (s *Scanner) scan() {
	// Go to each rule passing the scanner?
	// Create sub scanner for each ruller?
	//

}

//Extractors
// Scan word
func word(s *Scanner, word string) bool {

}
