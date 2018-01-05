package lexer

import (
	"fmt"
	"strings"
	"testing"
)

// Stack based rules
//
// Definitions
//   define what is what
//   KEYWORD: "select,update,test"
// root {
//   Define a command
//   {SELECT fields FROM table .. }
//
// }
// Define SELECT stmt A
//

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
		s := NewScanner(str)
		fmt.Println(str)
		for {
			var t synToken
			s.scan(&t)
			if t.id == EOF {
				break
			}
			fmt.Printf("%s^%s Token: %s (%d) - %s\n",
				strings.Repeat(" ", t.pos),
				strings.Repeat("-", mlen-t.pos),
				TokenStr[t.id], t.pos, t.str,
			)
			//s.Suggestion() ??
		}
	}

}
