package lexrp

import "fmt"

// Definitions actually
type ItemType int

type Item struct {
	Typ ItemType
	Val string
	Pos int // Extra position in the input string
}

func (i Item) String() string {
	switch i.Typ {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return i.Val
	}

	if len(i.Val) > 10 {
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}
