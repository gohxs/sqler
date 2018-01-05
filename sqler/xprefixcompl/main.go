// Just a test on string base completion
package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gohxs/termu"
)

var reSpace = regexp.MustCompile("\\s+")

var (
	entry = citem(
		citem("select",
			citem("from",
				citem(func() []string {
					now := time.Now()
					return []string{"db1", "cmwoms", "othertable", now.Format("05")} // dynamic
				}, "maintable"),
			),
			citem("from2"),
			"insert",
			"update",
			"delete"),
	)
)

func main() {

	t := termu.New()
	t.SetPrompt("prefix> ")
	t.AutoComplete = compl(t)
	log.SetOutput(t)

	for {
		line, err := t.ReadLine()
		if err != nil {
			break
		}
		fmt.Fprint(t, line)

	}
}

// Completion thing
func compl(t *termu.Term) func(line string, pos int, key rune) (string, int, bool) {
	return func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		if key != '\t' {
			return
		}

		res, rest := entry.search(reSpace.Split(line, -1))
		if len(res) == 0 {
			return "", 0, false // Nothing to do
		}
		if rest != "" { // Add the common match
			newLine = line + rest
			newPos = len(newLine)
			ok = true
		}
		if len(res) > 1 { // Print to user if more than 1 results
			fmt.Fprintln(t, res)
		}

		return

	}
}

// Recurse a tree and find next word?

// Search tree:
// root
//   word1
//		sub1
//		sub2
//			sub3
//   word2
//   word3

/*type prefixEntry interface {
	Names() []string
	Children() []prefixEntry
	search(parts []string) ([]string, string)
}*/

type SEntry struct {
	names func() []string // List of strings? function to get a list
	// Sub entries?
	children []*SEntry
}

func (s *SEntry) search(parts []string) ([]string, string) {
	//log.Println("Search for", parts)
	//log.Println("names:", s)
	ret := []string{}
	if len(parts) == 0 { // Return this one
		for _, v := range s.children {
			ret = append(ret, v.names()...) // All names
		}
		return ret, ""
	}
	// If our name is exacly part[0[
	for _, v := range s.children {
		for _, n := range v.names() {
			if strings.HasPrefix(n, parts[0]) {
				ret = append(ret, n)
			}
			if len(parts) > 1 && n == parts[0] { // If exacly equal we search forward and last match
				//log.Printf("found at: '%s' next will be '%#v'", v.names(), parts[1:])
				return v.search(parts[1:])
			}
		}
	}
	// Completion rest
	rest := ""
	if len(ret) == 1 {
		rest = ret[0][len(parts[0]):] + " "
	} else if len(ret) > 1 {
		match := commonPrefix(ret)
		rest = match[len(parts[0]):]
	}
	// Match thing
	return ret, rest
}

func citem(params ...interface{}) *SEntry {
	ret := &SEntry{children: []*SEntry{}}
	nameR := []string{}
	fnList := []func() []string{}
	for _, p := range params {
		switch v := p.(type) {
		case string:
			nameR = append(nameR, v)
		case func() []string:
			fnList = append(fnList, v)
		case *SEntry:
			ret.children = append(ret.children, v)
		}
	}
	ret.names = func() []string {
		r := []string{}
		r = append(r, nameR...)
		for _, fn := range fnList {
			r = append(r, fn()...)
		}
		return r
	}

	return ret
}
