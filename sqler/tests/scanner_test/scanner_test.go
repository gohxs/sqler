package testgoscanner

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type testable struct {
	name string
	fn   func(string, tokenHandler)
}

var (
	/*src = `
	CREATE DATABASE "test";
	CREATE TABLE "user" ("id" SERIAL, "name" string, "phone" string);
	INSERT INTO "table" (id, name) values(1, '1');
							` /**/
	src = `INSERT INTO "table" (id, name, float, hex, octal) values(1,'12',3e2, 0xAFaF,0644)`

	/*src = `
	10.2, 0x100, 10e4, 10.2.1, 0644
	10.2, 0x100, 10e4, 10.2.1, 0644
	`*/

	testFuncs = []testable{
		{"pkg/text/scanner ", goTextScanner},
		{"pkg/go/scanner   ", goScanner},
		//{"rob pike", robPikeScanner},
		{"fatih/hcl/scanner", fatihScanner},
		{"sqler/hxs_struct ", hxsTextScanner_struct},
		{"sqler/hxs_state  ", hxsTextScanner_state},
	}
)

func TestAll(t *testing.T) {
	for _, ta := range testFuncs {
		t.Run(ta.name, func(t *testing.T) {
			fmt.Println("Func:", ta.name)
			fmt.Printf("%s\n", src)
			ta.fn(src, doPrintThings)
		})
	}
}

var benchSrc string

func init() {
	// Repeat source
	benchSrc = ""
	for i := 0; i < 1000; i++ { // Parse 1000 LOC
		benchSrc += src
	}
}

func BenchmarkNormal(b *testing.B) {
	benchHandler(b, loopBench)
}

func BenchmarkStuff(b *testing.B) {
	benchHandler(b, loopBenchStuff)
}
func BenchmarkDelay(b *testing.B) {
	benchHandler(b, loopBenchDelay)
}

func BenchmarkChannel(b *testing.B) {
	benchHandler(b, loopBenchChannel)
}

func BenchmarkChannelDelay(b *testing.B) {
	benchHandler(b, loopBenchDelay)
}

var mlen = 80

// Test all the bench in testFuncs
func benchHandler(b *testing.B, fn looper) {
	for _, ta := range testFuncs {
		b.Run(ta.name, func(b *testing.B) {
			// Looper
			fn(b, ta.fn)
		})
	}
}

type looper func(*testing.B, func(string, tokenHandler))

func loopBench(b *testing.B, fn func(string, tokenHandler)) {
	for i := 0; i < b.N; i++ {
		fn(benchSrc, doBenchDoer)
	}
}
func loopBenchStuff(b *testing.B, fn func(string, tokenHandler)) {
	for i := 0; i < b.N; i++ {
		fn(benchSrc, doBenchStuff)
	}
}

func loopBenchDelay(b *testing.B, fn func(string, tokenHandler)) {
	for i := 0; i < b.N; i++ {
		fn(benchSrc, doBenchDoerDelay)
	}
}
func loopBenchChannel(b *testing.B, fn func(string, tokenHandler)) {
	for i := 0; i < b.N; i++ {
		channelize(benchSrc, fn, doBenchDoer)
	}
}
func loopBenchChannelDelay(b *testing.B, fn func(string, tokenHandler)) {
	for i := 0; i < b.N; i++ {
		channelize(benchSrc, fn, doBenchDoerDelay)
	}
}

func color(n int) string {

	colorStart := 84
	colorLast := 231

	n *= 16
	n = n % (colorLast - colorStart) // limit to 100
	n += colorStart                  // Start on color 100
	return fmt.Sprintf("\033[38;5;%dm", n)
}

// Space to improve but its working
// Print each token, colored
func doPrintThings(tok *tokSample) {
	fmt.Printf("%s", color(tok.count))

	col := tok.col() - 1
	line := tok.line()
	fmt.Printf("%s %s  %-14.14s \033[01;35m%-14.14s (%2d) \033[0m at: \033[01;30m%2d,%-2d offset: %3d\033[0m",
		strings.Repeat(" ", col),
		strings.Repeat("\u2500", mlen-col), // Line
		fmt.Sprintf("%#v", tok.val()),
		tok.name, tok.ID,
		line, col, tok.off(),
	)
	fmt.Printf("%s", color(tok.count))

	fmt.Printf("\033[%dA\r", tok.count+1)
	if col > 0 {
		fmt.Printf("\033[%dC", col)
	}
	fmt.Printf("\033[7m%s\033[0m\n", tok.val()) // Value thing

	fmt.Printf("%s", color(tok.count))
	// For each line
	// We are upthere go down one by one printing the Vertical lines
	if col > 0 {
		fmt.Printf("\033[%dC", col)
	}
	fmt.Printf("\u25B3\n") // Arrow

	for i := 0; i < tok.count; i++ { // We are upd
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

func doBenchStuff(tok *tokSample) {

	// grab Line etc from each tok
	// Abuse grabbers not always we need this
	tok.val()
	tok.line()
	tok.col()
	tok.off()

}

// Token handlers
func doBenchDoerDelay(tok *tokSample) {
	doBenchDoer(tok)
	<-time.After(100) // Simulate delay, i.e. parser working
}

// Tiny handler
func doBenchDoer(tok *tokSample) {
	//if tok.ID == scanner.Ident {
	if strings.HasPrefix("INSERT", tok.val()) {
	}
	//}
}

// Handle token and send to chan
func doChannelize(tokenChan chan tokSample) func(*tokSample) {
	return func(tok *tokSample) {
		// We need to clone values because we are async
		text := tok.val()  // Libs dont work well with async
		col := tok.col()   // Libs dont work well with async
		line := tok.line() // Libs dont work well with async
		off := tok.off()   // Libs dont work well with async

		tok.val = func() string { return text }
		tok.col = func() int { return col }
		tok.line = func() int { return line }
		tok.col = func() int { return off }

		tokenChan <- *tok // Copy
	}
}

func channelize(src string, scanner func(string, tokenHandler), fn tokenHandler) {
	tokenChan := make(chan tokSample, 10)
	go func() {
		scanner(src, doChannelize(tokenChan))
		close(tokenChan)
	}()
	// get token from chan and call the handler
	for t := range tokenChan {
		fn(&t)
	}
}
