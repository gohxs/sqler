package testgoscanner

import (
	"fmt"
	goscanner "go/scanner"
	"go/token"
	"strings"
	"text/scanner"

	fscanner "github.com/fatih/hcl/scanner"
	ftoken "github.com/fatih/hcl/token"

	hscanner_state "github.com/gohxs/sqler/sqler/scanner"
	hscanner_struct "github.com/gohxs/sqler/sqler/tests/hxsscanner.struct"
)

// tokSample generic token taken from each lexer
// the reason there are funcs is that
// some lexers only do some processing if funcs are called
// so to better test a lexer we pass func wrappers
// that in most of cases it will only call the underlying func
// if called
type tokSample struct {
	ID    rune
	name  string
	count int
	val   func() string

	col  func() int
	line func() int
	off  func() int
}

// tokenWalker
type tokenHandler func(tok *tokSample)

// Per token
// ///////////////
//The following funcs will go through each token
//and pass a tokSample
func hxsTextScanner_state(src string, fn tokenHandler) {
	s := hscanner_state.New([]byte(src))
	for count := 0; ; count++ {
		tok := s.Scan()
		if tok == hscanner_state.EOF {
			break // done
		}
		fn(&tokSample{
			ID:    rune(tok),
			name:  tok.String(),
			count: count,
			val:   s.TokenText, // pass func
			off:   func() int { return s.Position.Offset },
			col:   func() int { return s.Position.Column },
			line:  func() int { return s.Position.Line },
		})
	}
}

func hxsTextScanner_struct(src string, fn tokenHandler) {
	s := hscanner_struct.New([]byte(src))
	for count := 0; ; count++ {
		if !s.Scan() {
			break
		}
		tok := s.Token()
		fn(&tokSample{
			ID:    rune(tok.Typ),
			name:  hscanner_struct.TokenStr[tok.Typ],
			count: count,
			val:   tok.Text, // func
			off:   func() int { return tok.Pos.Offset },
			col:   func() int { return tok.Pos.Column },
			line:  func() int { return tok.Pos.Line },
		})
	}
}

func goScanner(src string, fn tokenHandler) {
	s := goscanner.Scanner{}
	fset := token.NewFileSet()
	ffile := fset.AddFile("", fset.Base(), len(src))
	s.Init(ffile, []byte(src), func(p token.Position, msg string) {}, goscanner.Mode(0xFFFFFFFF)) // All

	for count := 0; ; count++ {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fn(&tokSample{
			ID:    rune(tok),
			name:  fmt.Sprintf("%s", tok),
			count: count,
			val:   func() string { return lit }, // pass func
			off:   func() int { return int(pos) - 1 },
			// Did not found exposed line Counters
			col: func() int {
				lastLine := strings.LastIndex(src[:pos], "\n") + 1 // current col
				return len(src[lastLine:pos])
			},
			line: func() int {
				return strings.Count(src[:pos], "\n") + 1 // Manual liner
			},
		})
	}

}
func goTextScanner(src string, fn tokenHandler) {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(src))
	s.Error = func(s *scanner.Scanner, msg string) {}
	for count := 0; ; count++ {
		tok := s.Scan()
		if tok == scanner.EOF {
			break
		}
		fn(&tokSample{
			ID:    tok,
			name:  scanner.TokenString(tok),
			count: count,
			val:   s.TokenText, // pass func
			off:   func() int { return s.Position.Offset },
			col:   func() int { return s.Position.Column },
			line:  func() int { return s.Position.Line },
		})
	}
}

func fatihScanner(src string, fn tokenHandler) {
	s := fscanner.New([]byte(src))
	for count := 0; ; count++ {
		tok := s.Scan()
		if tok.Type == ftoken.EOF {
			break
		}
		fn(&tokSample{
			ID:    rune(tok.Type),
			name:  fmt.Sprintf("%s", tok.Type),
			count: count,
			val:   func() string { return tok.Text },
			off:   func() int { return tok.Pos.Offset },
			col:   func() int { return tok.Pos.Column },
			line:  func() int { return tok.Pos.Line },
		})
	}
}
