package sqler

import (
	"bytes"
	"fmt"
	"io"
)

//Returns a common prefix in a list
func commonPrefix(list []string) string {
	if len(list) == 0 {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	for i, ch := range list[0] {
		for _, v := range list[1:] {
			if i >= len(v) || ch != []rune(v)[i] {
				return buf.String()
			}
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}

//FmtWriter Wraps funcs like Println,Print,Printf to a writer
type FmtWriter struct {
	io.Writer
}

//Println println implementation
func (f *FmtWriter) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(f, a...)
}

//Printf printf implementation
func (f *FmtWriter) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f, format, a...)
}

//Print print implementation
func (f *FmtWriter) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(f, a...)
}
