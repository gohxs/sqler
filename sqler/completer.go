package sqler

import (
	"fmt"
	"strings"

	"github.com/gohxs/termu"
)

type Suggestion struct {
	suggestion string
	stype      int
}

type ComplEngine struct {
	Suggest func(line string) ([]string, string)
	term    *termu.Term
	tab     int
	hist    int
	mode    int
	maxCol  int
	history []string
	suggest []string // suggestions, maybe not needed
}

///////////////////////
// Mode 0 - show history
// Mode 1 - show completion
func (c *ComplEngine) Display(in string) string {
	if len(in) == 0 { // pass right throuh its a 0
		return in
	}
	res := highlight(in)
	if len(c.history) > 0 && c.mode == 0 {
		m := c.history[c.hist%len(c.history)] // Select one from list
		if len(in) < len(m) {
			c := m[len(in):] // Find space
			if sp := strings.Index(c, " "); sp != -1 {
				c = c[:sp]
			}

			res += "\033[01;30m" + c + "\033[m"
		}
	}
	if len(c.suggest) == 0 {
		return res
	}
	if c.mode != 1 {
		return res
	}
	/////////////////////////////////////////
	// Show the suggestions as a menu
	////////////////////////////////////////////
	maxLen := 0
	for _, s := range c.suggest {
		if l := len(s); l > maxLen {
			maxLen = l
		}
	}
	maxLen += 2 // pad
	width, _ := c.term.GetSize()
	c.maxCol = width / maxLen
	res += "\n" // Last word

	for i, s := range c.suggest {

		if i%c.maxCol == 0 && i != 0 {
			res += "\n" // move line
		}
		if i == c.tab%len(c.suggest) {
			res += fmt.Sprintf("\033[01;47m%-[1]*s\033[m", maxLen, s)
			continue
		}
		res += fmt.Sprintf("%-[1]*s", maxLen, s)
	}
	return res + "\033[J"
}

// AutoComplete compatible with golang.org/x/crypto/ssh/terminal
// Per key completion
func (c *ComplEngine) AutoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	//log.Printf("Running key: %#v %#v", key, '\n')
	switch key {
	case '\r', '\n':
		// Nothing to do
		if c.mode != 1 {
			return
		}
		c.mode = 0
		if len(c.suggest) == 0 {
			return
		}
		s := c.suggest[c.tab%len(c.suggest)]
		li := strings.LastIndex(line, " ") + 1
		if strings.HasPrefix(s, line[li:]) {
			newLine = line + s[len(line[li:]):] + " "
		} else {
			newLine = line + " " + s + " "
		}
		newPos = len(newLine)
		ok = true
		c.history = histMatchList(c.term, newLine) // On type
		return
	case termu.AKBackspace:
		c.mode = 0
		c.history = []string{} // clear history
	case termu.AKHistNext:
		//log.Println("Hist NexT")
		if c.mode == 1 {
			c.tab += c.maxCol
			if c.tab >= len(c.suggest) {
				c.tab = c.tab % c.maxCol
			}
			return line, pos, true // consume key
		}
	case termu.AKHistPrev:
		//log.Println("Hist NexT")
		if c.mode == 1 {
			c.tab -= c.maxCol
			if c.tab < 0 {
				c.tab += len(c.suggest) + (c.maxCol - len(c.suggest)%c.maxCol)
				if c.tab >= len(c.suggest) {
					c.tab -= c.maxCol
				}
			}
			return line, pos, true
		}
		c.history = nil
	case termu.AKCursLeft:
		if c.mode == 1 {
			c.tab--
			if c.tab < 0 {
				c.tab += len(c.suggest)
			}
			return line, pos, true
		}
	case termu.AKCursRight:
		if c.mode == 1 { // Menu auto complete
			c.tab++
			if c.tab >= len(c.suggest) {
				c.tab = c.tab % len(c.suggest)
			}
			return line, pos, true
		}
		if pos == len(line) {
			return c.histComplete(line, pos, false)
		}
		// Complete history
		/*if len(c.history) > 0 {
			newLine = c.history[c.hist%len(c.history)]
			return newLine, len(newLine), true
		}*/

	case termu.AKShiftTab:
		if c.mode == 1 { // next tab
			c.tab--
			if c.tab < 0 {
				c.tab = len(c.suggest) - 1
			}
		}
		return
	case '\t': // Tab
		c.history = histMatchList(c.term, line) // On type
		if c.mode == 1 {                        // next tab
			c.tab++
			if c.tab > len(c.suggest) {
				c.tab = c.tab % len(c.suggest)
			}
		}

		newLine, newPos, ok = c.tabComplete(line, pos)
		if !ok {
			newLine, newPos, ok = c.histComplete(line, pos, true)
		}
		return
		// Let the code go
	default:
		//While typing we load history (might be slow)
		//
		c.history = histMatchList(c.term, line+string(key))
		// If we are in menu mode, we reload the suggestions
		// Until no suggestions available
		if c.mode == 1 {
			c.suggest, _ = c.Suggest(line + string(key))
			if len(c.suggest) < 1 {
				c.mode = 0
			}
		}
		c.tab = 0
		return
	}
	return "", 0, false
	// Tab key
}

func (c *ComplEngine) tabComplete(line string, pos int) (newLine string, newPos int, ok bool) {

	//c.mode = 0 // Reset mode?

	li := strings.LastIndex(line, " ") + 1
	var compl string
	c.suggest, compl = c.Suggest(line)
	if len(c.suggest) > 1 {
		c.mode = 1
		newLine = line + compl
		return newLine, len(newLine), true
	}

	// Single match complete
	if len(c.suggest) == 1 {
		res := c.suggest[0]
		if strings.HasPrefix(res, line[li:]) {
			newLine = line[:li] + res + " "
		} else {
			newLine = line + " " + res + " "
		}
		return newLine, len(newLine), true
	}
	return "", 0, false
}

// histcomplete
// search == false, it will complete current history
// search == true, if history is more than 1 it will return change the index
func (c *ComplEngine) histComplete(line string, pos int, search bool) (newLine string, newPos int, ok bool) {

	if len(c.history) == 0 { // Fill the history
		//	c.history = histMatchList(c.term, line)
		return
	}

	var res string
	if len(c.history) > 1 && search == true {
		res = commonPrefix(c.history)
		//log.Printf("Get the common prefix to complete %#v", res)
		subHist := histMatchList(c.term, res)
		//log.Println("Sub history:", subHist)
		if len(subHist) > 0 {
			c.hist++
		}
	} else {
		// fetch the current one
		res = c.history[c.hist%len(c.history)]
	}

	// Skip any blanks Trim counting the offset
	for pos < len(res) && res[pos] == ' ' {
		pos++
	}

	space := 0
	if pos < len(res) { // Find next history spacespace
		space = strings.Index(res[pos:], " ") // current position forward
	}

	if space > 0 { // only positive space
		res = res[:pos+space] + " "
	} else {
		//if res[len(res)-1] != ' ' {
		//	res = res + " "
		//}
	}
	if res != line { // reset tab if line changed
		c.tab = 0
	}
	// Go to next space only
	// End complete
	return res, len(res), true

}

// Uniquelly List history
func histMatchList(t *termu.Term, in string) []string {
	if len(in) == 0 {
		return []string{}
	}
	exists := map[string]bool{}

	ret := []string{}

	histList := t.History.List()
	for i := len(histList) - 1; i >= 0; i-- {
		v := histList[i] // reverse
		if _, ok := exists[v]; ok {
			continue //already exists
		}
		if strings.HasPrefix(v, in) { // Print in in white, rest in black
			exists[v] = true
			ret = append(ret, v)
		}
	}
	return ret
}
