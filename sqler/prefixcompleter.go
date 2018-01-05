package sqler

import (
	"regexp"
	"strings"
)

type SEntry struct {
	names func(in string) []string // List of strings? function to get a list
	// Sub entries?
	children []*SEntry
}

func (s *SEntry) search(parts []string) ([]string, string) {
	ret := []string{}

	if len(parts) == 0 { // Return this one
		for _, v := range s.children {
			ret = append(ret, v.names("")...) // All names
		}
		return ret, commonPrefix(ret) // less len parts?
	}
	word := parts[0]
	for _, v := range s.children {
		for _, n := range v.names(word) {
			if len(parts) >= 1 && n == word { // If exacly equal we search forward and last match
				return v.search(parts[1:])
			}
			if strings.HasPrefix(n, word) {
				ret = append(ret, n) // Word parts SWITCH THIS HERE
			}
		}
	}
	var compl string
	if len(ret) == 1 {
		compl = ret[0][len(word):]
	}
	// Match thing
	return ret, compl
}

func citem(params ...interface{}) *SEntry {
	ret := &SEntry{children: []*SEntry{}}
	nameR := []string{}
	fnList := []func(string) []string{}
	for _, p := range params {
		switch v := p.(type) {
		case string:
			nameR = append(nameR, v)
		case func(string) []string:
			fnList = append(fnList, v)
		case func() []string:
			fnList = append(fnList, func(string) []string { return v() }) // wrapper
		case *SEntry:
			ret.children = append(ret.children, v)
		}
	}
	ret.names = func(word string) []string {
		r := []string{}
		r = append(r, nameR...)
		for _, fn := range fnList {
			r = append(r, fn(word)...)
		}
		return r
	}

	return ret
}

var reSpace = regexp.MustCompile("\\s+")
