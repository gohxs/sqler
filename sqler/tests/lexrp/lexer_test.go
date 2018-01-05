package lexrp

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {

	stmt := []string{
		`SEL`, // Should return suggestions for select
		`CREATE TABLE "table" (id varchar)`,
		`SELECT * FROM table`,
		`UPDATE "table" SET name = '3'`,
		`UPDATE "table" SET name = '3' WHERE`,
	}

	mlen := 40
	for _, str := range stmt {
		fmt.Println(str)

		_, items := lex("t", str)
		for it := range items {
			fmt.Printf("%s^%s Token: %s  %#v:%d\n",
				strings.Repeat(" ", it.pos),
				strings.Repeat("-", mlen-it.pos),
				TokenStr[it.typ],
				it.val, it.pos,
			)

			if it.typ == itemEOF {
				break
			}
		}
		//s.Suggestion() ??
	}

}
