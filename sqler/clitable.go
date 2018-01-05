package sqler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
)

// DataTable contains data to print
type DataTable [][]interface{}

func FromInterface(o interface{}) DataTable {
	d := [][]interface{}{}
	val := reflect.ValueOf(o)

	if val.Kind() != reflect.Slice {
		log.Println("O is not a slice")
		return nil // Empty
	}
	// Each element somehow

	el := reflect.ValueOf(o).Index(0)
	header := []interface{}{}
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		header = append(header, f.Type().Name())
	}
	d = append(d, header)

	// Elements
	for i := 0; i < val.Len(); i++ {
		row := []interface{}{}
		el := val.Index(i)
		for i := 0; i < el.NumField(); i++ {
			f := el.Field(i)
			row = append(row, fmt.Sprintf("%v", f.Interface()))
		}
		d = append(d, row)
	}
	return d

}

// Sprint to string: fmt.Sprint
func TableSprint(dataTable DataTable) string {
	var colMaxWidth []int
	if len(dataTable) == 0 {
		return "(no data)\n"
	}

	// Calc colWidth
	colMaxWidth = make([]int, len(dataTable[0]))

	for _, dataRow := range dataTable {
		for i, c := range dataRow {
			slen := len(c.(string))
			if colMaxWidth[i] < slen {
				colMaxWidth[i] = slen
			}
		}
	}

	headerElem := "|--"
	var formatElem []string
	for _, v := range colMaxWidth {
		formatElem = append(formatElem, fmt.Sprintf("  %%%ds", v))
		headerElem += strings.Repeat("-", v+4) + "|"
	}
	// The row with formats
	formatString := "|  " + strings.Join(formatElem, "  |") + "  |"

	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "\r\n")

	for i, dataRow := range dataTable {
		fmt.Fprintf(buf, formatString, dataRow...)
		if i == 0 { // Header line
			fmt.Fprintf(buf, "\r\n%s", headerElem)
		}

		/*str = fmt.Sprintf(formatString, dataRow...)
		if i == 0 {
			str += "\r\n"
			str += headerElem

			//+ fmt.Sprintf("|"+strings.Repeat("-", len(str)-2)+"|")
		}*/
		fmt.Fprint(buf, "\r\n")
	}
	fmt.Fprintf(buf, "\r\n(%d rows)\r\n", len(dataTable)-1)
	fmt.Fprintf(buf, "\r\n")

	return buf.String()
}

// Fprint print table to writer
func TableFprint(w io.Writer, dataTable DataTable) (int, error) {
	return fmt.Fprint(w, TableSprint(dataTable))
}

// Print print using fmt.Print
func TablePrint(dataTable DataTable) (int, error) {
	return fmt.Print(TableSprint(dataTable))
}
