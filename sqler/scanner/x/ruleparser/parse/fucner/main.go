package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gohxs/prettylog"
)

func main() {
	prettylog.Global()
	s := &state{}

	// build
	fn := Seq("string", "string2")
	// Execute
	fn(s)
	log.Println(s)
}

type state struct {
	data []string
}

func (s *state) add(a ...interface{}) {
	s.data = append(s.data, fmt.Sprint(a...))
}
func (s *state) addf(format string, a ...interface{}) {
	s.data = append(s.data, fmt.Sprintf(format, a...))
}
func (s *state) String() string {
	return strings.Join(s.data, "\n")
}

// Return a single func passing the head thing
// The recurseive func
type walkerFn func(s *state) walkerFn

func Seq(args ...interface{}) walkerFn {

	return func(s *state) walkerFn {
		for _, a := range args {
			ftyp(a)(s)
		}
		return nil
	}
}

func ftyp(t interface{}) walkerFn {

	switch v := t.(type) {
	case string:
		return func(s *state) walkerFn {
			s.add(v)
			return nil // no more or next?
		}
	}
	return nil
}
