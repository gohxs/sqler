package scannerutils

import (
	"fmt"
	"strings"
)

func color(n int) string {

	colorStart := 84
	colorLast := 231

	n *= 16
	n = n % (colorLast - colorStart) // limit to 100
	n += colorStart                  // Start on color 100
	return fmt.Sprintf("\033[38;5;%dm", n)
}

var mlen = 80

// Space to improve but its working
// Print each token, colored
func Print(ID int, name string, line, col int, text string, count int) {
	fmt.Printf("%s", color(count))

	col-- // remove 1

	fmt.Printf("%s %s  %-14.14s \033[01;35m%-14.14s (%2d) \033[0m at: \033[01;30m%2d,%-2d\033[0m",
		strings.Repeat(" ", col),
		strings.Repeat("\u2500", mlen-col), // Line
		fmt.Sprintf("%#v", text), name, ID,
		line, col,
	)
	fmt.Printf("%s", color(count))

	fmt.Printf("\033[%dA\r", count+1)
	if col > 0 {
		fmt.Printf("\033[%dC", col)
	}
	fmt.Printf("\033[7m%s\033[0m\n", text) // Value thing

	fmt.Printf("%s", color(count))
	// For each line
	// We are upthere go down one by one printing the Vertical lines
	if col > 0 {
		fmt.Printf("\033[%dC", col)
	}
	fmt.Printf("\u25B3\n") // Arrow

	for i := 0; i < count; i++ { // We are upd
		if col > 0 {
			fmt.Printf("\033[%dC", col)
		}
		fmt.Printf("\u2502\n") // Vertical line
	}
	fmt.Printf("\033[A")
	if col > 0 {
		fmt.Printf("\033[%dC", col)
	}
	fmt.Printf("\u2570\n") // corner

	fmt.Print("\033[m")

}
