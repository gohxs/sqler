package sqler_test

import (
	"testing"

	"github.com/gohxs/sqler/sqler"
)

type User struct {
	ID   int
	Name string
	Test int
}

func TestVisual(t *testing.T) {

	userList := []User{
		{ID: 1, Name: "user1", Test: 10},
		{ID: 2, Name: "this is suposed to be a huge field", Test: 30},
		{ID: 3, Name: "user3", Test: 40},
		{ID: 4, Name: "user4", Test: 10},
	}

	sqler.TablePrint(sqler.FromInterface(userList))

}
