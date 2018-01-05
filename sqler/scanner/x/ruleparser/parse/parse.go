package parse

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gohxs/prettylog"
	"github.com/gohxs/sqler/sqler/scanner"
	"github.com/gohxs/sqler/sqler/scanner/scannerutils"
)

type token struct {
	tok    scanner.TokenType
	text   string
	Line   int
	Column int
}

// Not sure if ideal to get a list of tokens right away
// Read about AST's to find a way to walk on the tree
// Or should this be an AST?
func tokenList(s *scanner.Scanner) []*token {
	tokens := []*token{}
	for {
		tok := s.Scan()
		if tok == scanner.EOF {
			break
		}
		tokens = append(tokens, &token{
			tok:    tok,
			text:   s.TokenText(),
			Line:   s.Position.Line,
			Column: s.Position.Column,
		})
	}
	return tokens
}
func tokenIter(s *scanner.Scanner) *Iter {
	tokens := tokenList(s)
	return &Iter{tokens, -1}
}

func tokenParse(src string) *Iter {
	return tokenIter(scanner.New([]byte(src)))
}

// Test rules
//
func Start() {
	// Testing the iterator
	//
	walkers := []walkerFn{}

	walkers = append(walkers, fnFromType("string"))

	prettylog.Global()
	src := "SELECT * FROM WHATEVER"
	it := tokenIter(scanner.New([]byte(src)))
	fmt.Println(src)

	for count := 0; ; count++ {
		t := it.Next()
		if t == nil {
			break
		}
		scannerutils.Print(int(t.tok), t.tok.String(), t.Line, t.Column, t.text, count)
	}
	// Define rules somehow
	log.Println("Build the sample RULE")
	runner := Root(
		Rule("SELECT", Join(scanner.IDENT, ","), "FROM", Tables),
		//Rule("INSERT", "INTO", Tables, Opt(Rule(scanner.LPAREN, OneOrMore(Rule(scanner.IDENT, Opt(scanner.RPAREN))))), "VALUES"),
	)

	log.Println("Run the walker rules")
	runner(it.Clone().Reset())
	//runner(tokenIter(scanner.New([]byte("S")))) // Give SUGGESTIONS about words starting with S

}

// Returns a list of tables
func Tables() []string {
	return []string{"test", "test2"}
}

// First pass Type to walkerFn with specific fns
// Create a fn Iterator that calls func to func
// Go through the fns and check the it result

// The function that understands the scanner/Iterator
//type walkerFn func(it *Iter) walkerFn
type walkerFn func(it *Iter) walkerFn

func Root(rules ...walkerFn) walkerFn {
	return func(it *Iter) walkerFn {
		newIt := it.Clone() // New iterator from start
		// Each rule?
		for _, r := range rules {
			r(newIt)
		}
		return nil
	}
}

// Rule create a sequential aware walkerFn
func Rule(a ...interface{}) walkerFn {
	log.Println("Creating a rule func")
	// Create walker arrayhere

	fnList := []walkerFn{}
	for i, arg := range a {
		log.Printf("Get walkerfn For Rule #%d - %s", i, reflect.TypeOf(a).Elem().Name())
		fnList = append(fnList, fnFromType(arg))
	}

	return func(it *Iter) walkerFn {
		log.Println("Executing a Rule")
		log.Println("Cloning iterator")
		// Sequencializer instead of for loop
		nit := it.Clone()
		for i, fn := range fnList {
			log.Printf("Rule Seq: #%d", i)
			if fn == nil { // Ignore?
				return nil
			}
			fn(nit)
			nit.Next()
			//Search case?
		}
		return nil
	}
}

// Join, join one or more with a key
//   -------> Rule("SELECT",Join(Or(IDENT,"*"), ","))
//   ------> Perfect join match
func Join(m interface{}, sep interface{}) walkerFn {
	log.Println("Preparing join matcher with ", reflect.TypeOf(m).Name())
	expectFn := fnFromType(m)
	sepFn := fnFromType(sep)
	return func(it *Iter) walkerFn {
		r := expectFn(it)
		if r != nil {
			// Not sure what to do

		}
		it.Next()
		r = sepFn(it) // Should we move next here or before?
		if r == nil {
			// Not sure what to do

		}
		log.Printf("Matching the join func Expecting: %v, separated by : %v", m, sep)

		return nil
	}
}

// Create a walker Function from a type
func fnFromType(a interface{}) walkerFn {
	switch v := a.(type) {
	case string: // Expect exact key word
		log.Printf("Arg is \"%s\", returning string matcher walker", v)
		return func(it *Iter) walkerFn { // The runner func
			t := it.Peek()
			if t == nil {
				return nil
			}
			log.Printf("Expecting key word: \"%s\", found: \"%s\"", v, t.text)
			tex := t.text
			if tex != v {
				// add error / suggestion to iter
				_ = v
			}
			// Return next?
			return nil
		}
	case scanner.TokenType:
		log.Printf("Arg is \"%s\", returning token matcher", v.String())
		return func(it *Iter) walkerFn {
			// Should be token type
			t := it.Peek()
			log.Printf("Matching with IDENT: %s, %s ", v.String(), t.tok)
			if t.tok != v {
				return nil
			}

			return nil
		}
	case walkerFn:
		return v
	}
	return nil

}

//Any of the A
func Or(a ...interface{}) walkerFn {

	return nil
}
func Opt(a ...interface{}) walkerFn {
	return func(*Iter) walkerFn {
		log.Println("Optionally")
		return nil
	}
}
func OneOrMore(a ...interface{}) walkerFn {
	return func(*Iter) walkerFn {
		log.Println("One or more")
		return nil
	}
}
