package sqler_test

import (
	"testing"

	"github.com/gohxs/sqler/sqler"
)

func TestPkg(t *testing.T) {

	s := sqler.New()

	s.Cmd(".open sqlite3 :memory:")
	s.Cmd("CREATE table hi (id int, man string)")
}
