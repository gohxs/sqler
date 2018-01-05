package sqler

import (
	"regexp"
	"strings"
)

func parseargs(line string) []string {
	re := regexp.MustCompile("\\s+")
	// Cool trick
	strPart := strings.Split(line, "\"")
	var args []string
	for i, v := range strPart {
		if i&0x1 == 0 {
			lparts := re.Split(strings.TrimSpace(v), -1)
			args = append(args, lparts...) // Add all parts if not inside '"'
		} else {
			args = append(args, v)
		}
	}

	return args
}
